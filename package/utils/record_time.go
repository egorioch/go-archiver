package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// CreateDataYMD Создание директории, согласно текущей дате
func CreateDataYMD(videosPath string, cameraID int) (string, error) {
	now := time.Now()

	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	outputDir := fmt.Sprintf("%s/%s/%s/%s/%d", videosPath, year, month, day, cameraID)

	// Создаем директорию, если её нет
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка при создании директории:", err)
		return "", err
	}

	return outputDir, nil
}

func GetShortURL(cameraURL string) string {
	return getStringInBetween(cameraURL, "//", ":")
}

// GetStringInBetween Returns empty string if no start string found
func getStringInBetween(str string, start string, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	e = e + s
	return str[s:e]
}
