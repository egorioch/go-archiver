package logger

import (
	"log"
	"os"
)

func CreateLogger(filePath string) (*log.Logger, error) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger, nil
}
