package folder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"roumet/logger"
)

var (
	// ErrNoImagesFound is returned when no images are found in the folder.
	ErrNoImagesFound = errors.New("aucune image trouvée dans le dossier")
)

func GetCurrentFolderName() (string, error) {
	// Récupère le chemin absolu du dossier de travail actuel
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Extrait juste le nom du dernier dossier
	return filepath.Base(workingDir), nil
}

func ReadFolderFilName(path string) ([]string, error) {

	var fileNames []string
	files, err := os.ReadDir(path)
	if err != nil {
		logger.Error("folder", err)
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		// Vérifie si le nom du fichier se termine par .jpg, .jpeg ou .png
		if fileName == "" || len(fileName) < 5 {
			continue
		}

		lastElement := fileName[len(fileName)-4:]
		if lastElement == ".jpg" || lastElement == ".png" || lastElement == ".jpeg" ||
			lastElement == ".JPG" || lastElement == ".PNG" || lastElement == ".JPEG" {
			fileNames = append(fileNames, fileName)
		}
	}

	if len(fileNames) == 0 {
		err := fmt.Errorf("floder: aucune image trouvée dans le dossier %s %w", path, ErrNoImagesFound)
		logger.Error("folder", err)
		return nil, ErrNoImagesFound
	}

	return fileNames, nil
}

// sort the file names
func SortImageFilenames(files []string) error {
	if len(files) == 0 {
		return errors.New("le tableau de fichiers est vide")
	}
	// On définit un regex pour extraire les numéros principaux et secondaires
	re := regexp.MustCompile(`^(\d+)(?:-(\d+))?\.jpg$`)

	sort.Slice(files, func(i, j int) bool {
		matchI := re.FindStringSubmatch(files[i])
		matchJ := re.FindStringSubmatch(files[j])

		if matchI == nil || matchJ == nil {
			// Si le nom ne matche pas le format attendu, on compare simplement par string
			return files[i] < files[j]
		}

		// Conversion des numéros en entiers
		mainI, _ := strconv.Atoi(matchI[1])
		mainJ, _ := strconv.Atoi(matchJ[1])

		if mainI != mainJ {
			return mainI < mainJ
		}

		// Comparaison des sous-numéros (si présents)
		var subI, subJ int
		if matchI[2] != "" {
			subI, _ = strconv.Atoi(matchI[2])
		}
		if matchJ[2] != "" {
			subJ, _ = strconv.Atoi(matchJ[2])
		}
		return subI < subJ
	})
	return nil
}

type FileInfoWithPath struct {
	Name     string
	FullPath string
	ModTime  time.Time
}

// ListFilesSortedByCreation retourne la liste des noms de fichiers d’un dossier, triés par date de modification
func ListFilesSortedByCreation(dir string) ([]string, error) {
	var files []FileInfoWithPath

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		files = append(files, FileInfoWithPath{
			Name:     info.Name(),
			FullPath: path,
			ModTime:  info.ModTime(),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Tri par date de modification croissante
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	})

	var result []string
	for _, f := range files {
		result = append(result, f.Name)
	}
	return result, nil
}
