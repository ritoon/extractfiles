package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	URL_ROUMET = "https://www.roumet.com/photos/"
)

func main() {
	logAction("--------------------------------------------------")
	logAction("Début du traitement")
	folerName, err := getCurrentFolderName()
	if err != nil {
		logError(err)
		return
	}
	logAction("Le nom du dossier actuel est : " + folerName)

	files, err := readFolderFilName("./")
	if err != nil {
		logError(err)
		return
	}
	logAction("Le nombre de fichiers dans le dossier est : " + strconv.Itoa(len(files)))

	err = sortImageFilenames(files)
	if err != nil {
		logError(err)
		return
	}

	logAction("Le tri des fichiers a été effectué avec succès")

	mapFiles := make(map[int]string)
	increment := 1
	for i := 0; i < len(files); i++ {
		if !strings.Contains(files[i], "-") {
			increment++
			mapFiles[increment] = URL_ROUMET + folerName + "/" + files[i]
		} else {
			mapFiles[increment] += "|" + URL_ROUMET + folerName + "/" + files[i]
		}
	}

	logAction("Le nombre de lots dans le dossier est : " + strconv.Itoa(len(mapFiles)))

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logError(err)
		}
	}()

	logAction("Création du fichier Excel avec succès")

	// Create a new sheet.
	index, err := f.NewSheet("images")
	if err != nil {
		logError(err)
		return
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	f.DeleteSheet("Sheet1")

	logAction("Création de la feuille Excel images avec succès")

	// f.DeleteSheet("Sheet1")
	f.SetColWidth("images", "B", "B", 100)
	setHeader(f, "images")

	logAction("Ajout des entêtes dans la feuille Excel images avec succès")

	for k, filepathname := range mapFiles {
		AKeystring := fmt.Sprintf("A%d", k)
		BKeystring := fmt.Sprintf("B%d", k)
		CKeystring := fmt.Sprintf("C%d", k)
		f.SetCellValue("images", AKeystring, extractLotNumber(filepathname))
		f.SetCellValue("images", BKeystring, filepathname)
		if strings.Contains(filepathname, "|") {
			elems := strings.Split(filepathname, "|")
			// nbImg := strings.Count(filepathname, "|")
			f.SetCellValue("images", CKeystring, len(elems))
		} else {
			f.SetCellValue("images", CKeystring, 1)
		}
	}

	logAction("Ajout des données dans la feuille Excel images avec succès")

	// Save spreadsheet by the given path.
	if err := f.SaveAs("roumet-images.xlsx"); err != nil {
		logError(err)
	}
	logAction("Enregistrement du fichier Excel avec succès")
	logAction("Fin du traitement")
	logAction("--------------------------------------------------")
}

func getCurrentFolderName() (string, error) {
	// Récupère le chemin absolu du dossier de travail actuel
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Extrait juste le nom du dernier dossier
	return filepath.Base(workingDir), nil
}

func readFolderFilName(path string) ([]string, error) {

	var fileNames []string
	files, err := os.ReadDir(path)
	if err != nil {
		logError(err)
		return nil, err
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".jpg") || strings.Contains(file.Name(), ".JPG") ||
			strings.Contains(file.Name(), ".jpeg") || strings.Contains(file.Name(), ".JPEG") ||
			strings.Contains(file.Name(), ".png") || strings.Contains(file.Name(), ".PNG") {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil

}

// sort the file names
func sortImageFilenames(files []string) error {
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

func extractLotNumber(url string) string {

	// Le regex capture le nombre avant ".jpg" ou "-X.jpg"
	re := regexp.MustCompile(`/(\d+)(?:-\d+)?\.(?:jpg|jpeg|png)$`)
	match := re.FindStringSubmatch(url)
	if match == nil {
		logError(fmt.Errorf("Aucun numéro de lot trouvé dans l'URL : %s", url))
		return ""
	}
	return match[1]
}

var entetes = map[int]string{
	0: "Numéro de lot",
	1: "Images",
	2: "NB d'images",
}

func setHeader(f *excelize.File, sheetName string) {
	// Set header
	for i, v := range entetes {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellValue(sheetName, cell, v)
	}
}

func logError(err error) {
	if err == nil {
		return
	}

	f, fileErr := os.OpenFile("error.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Printf("Impossible d'ouvrir error.txt : %v", fileErr)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("Erreur : %v\n", err)
}

func logAction(msg string) {
	if msg == "" {
		return
	}

	f, fileErr := os.OpenFile("actions.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Printf("Impossible d'ouvrir error.txt : %v", fileErr)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("Erreur : %v\n", msg)
}
