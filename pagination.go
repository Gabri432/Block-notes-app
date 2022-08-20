package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

const NumPostFields = 5

func ReadFromTo(w http.ResponseWriter, fileName string, endPosition, startPosition int) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return []byte{}
	}
	defer file.Close()
	newContent := make([]byte, endPosition)
	_, errorMessage := file.ReadAt(newContent, int64(startPosition))
	if errorMessage != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return []byte{}
	}
	return newContent
}

func ScanPosts(file *os.File, startingLine, endingLine int) ([]byte, error) {
	if startingLine < 1 {
		startingLine = 1
	}
	if difference := endingLine - startingLine; difference%NumPostFields != 0 {
		endingLine = difference
	}
	currentLine := 0
	scanner := bufio.NewScanner(file)
	text := ""
	for scanner.Scan() {
		if currentLine < startingLine {
			currentLine++
			continue
		}
		if currentLine+1 > endingLine {
			break
		}
		text += scanner.Text()
	}
	bytes, err := json.MarshalIndent(text[:len(text)-1], "", " ")
	if err != nil {
		return []byte{}, err
	}
	return bytes, err
}
