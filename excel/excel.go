package excel

import (
	"fmt"
	"roumet/logger"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"roumet/folder"
	"roumet/util"
)

const (
	URL_ROUMET = "https://www.roumet.com/photos/"
)

func CreateExcelFile() error {
	logger.Action("--------------------------------------------------")
	logger.Action("Début du traitement")
	folerName, err := folder.GetCurrentFolderName()
	if err != nil {
		logger.Error("excel", err)
		return nil
	}
	logger.Action("Le nom du dossier actuel est : " + folerName)

	files, err := folder.ReadFolderFilName("./")
	if err != nil {
		logger.Error("excel", err)
		return nil
	}
	logger.Action("Le nombre de fichiers dans le dossier est : " + strconv.Itoa(len(files)))

	err = folder.SortImageFilenames(files)
	if err != nil {
		logger.Error("excel", err)
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

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("excel", err)
		}
	}()

	logger.Action("Création du fichier Excel avec succès")

	// Create a new sheet.
	index, err := f.NewSheet("images")
	if err != nil {
		logger.Error("excel", err)
		return nil
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	f.DeleteSheet("Sheet1")

	logger.Action("Création de la feuille Excel images avec succès")

	// f.DeleteSheet("Sheet1")
	f.SetColWidth("images", "B", "B", 100)
	setHeader(f, "images")

	logger.Action("Ajout des entêtes dans la feuille Excel images avec succès")

	for k, filepathname := range mapFiles {
		AKeystring := fmt.Sprintf("A%d", k)
		BKeystring := fmt.Sprintf("B%d", k)
		CKeystring := fmt.Sprintf("C%d", k)
		f.SetCellValue("images", AKeystring, util.ExtractLotNumber(filepathname))
		f.SetCellValue("images", BKeystring, filepathname)
		if strings.Contains(filepathname, "|") {
			elems := strings.Split(filepathname, "|")
			// nbImg := strings.Count(filepathname, "|")
			f.SetCellValue("images", CKeystring, len(elems))
		} else {
			f.SetCellValue("images", CKeystring, 1)
		}
	}

	logger.Action("Ajout des données dans la feuille Excel images avec succès")

	// Save spreadsheet by the given path.
	if err := f.SaveAs("roumet-images.xlsx"); err != nil {
		logger.Error("excel", err)
		return nil
	}
	logger.Action("Enregistrement du fichier Excel avec succès")
	logger.Action("Fin du traitement")
	logger.Action("--------------------------------------------------")
	return nil
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
