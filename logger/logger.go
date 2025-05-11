package logger

import (
	"log"
	"os"
)

const (
	ErrorFileName   = "errors.txt"
	ActionsFileName = "actions.txt"
)

func Error(lib string, err error) {
	if err == nil {
		return
	}

	f, fileErr := os.OpenFile(ErrorFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Printf("Impossible d'ouvrir error.txt : %v", fileErr)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("Erreur : %v %v\n", lib, err)
}

func Action(msg string) {
	if msg == "" {
		return
	}

	f, fileErr := os.OpenFile(ActionsFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Printf("Impossible d'ouvrir error.txt : %v", fileErr)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("Erreur : %v\n", msg)
}
