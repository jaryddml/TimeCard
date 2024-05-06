package main

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func getFileInfo() (string, string, error) {
	fileName := ""
	directoryPath := ""

	fileName = os.Args[0]

	directoryPath, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	return fileName, directoryPath, nil
}

func getFileDiscoveryDate(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}

func countLines(filePath string) (int, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	lines := strings.Count(string(content), "\n")
	return lines, nil
}

func getTimeSpentInFile(filePath string) (time.Duration, error) {
	return time.Duration(0), nil
}
