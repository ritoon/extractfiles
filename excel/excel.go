package excel

import (
	"fmt"
	"roumet/logger"
	"strings"

	"github.com/xuri/excelize/v2"

	"roumet/util"
)

func CreateExcelFile(mapFiles map[int]string) error {

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
