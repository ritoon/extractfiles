package main

import (
	"roumet/folder"
	"roumet/logger"
	"strconv"
	"strings"
)

const (
	URL_ROUMET = "https://www.roumet.com/photos/"
)

func main() {

	err := Run()
	if err != nil {
		logger.Error("main", err)
		return
	}
}

func Run() error {
	logger.Action("--------------------------------------------------")
	logger.Action("Début du traitement")
	folerName, err := folder.GetCurrentFolderName()
	if err != nil {
		logger.Error("main", err)
		return nil
	}
	logger.Action("Le nom du dossier actuel est : " + folerName)

	files, err := folder.ReadFolderFilName("./")
	if err != nil {
		logger.Error("main", err)
		return nil
	}
	logger.Action("Le nombre de fichiers dans le dossier est : " + strconv.Itoa(len(files)))

	err = folder.SortImageFilenames(files)
	if err != nil {
		logger.Error("main", err)
		return nil
	}

	logger.Action("Le tri des fichiers a été effectué avec succès")

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

	logger.Action("Le nombre de lots dans le dossier est : " + strconv.Itoa(len(mapFiles)))
	return nil
}
