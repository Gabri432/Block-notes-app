package main

import (
	"bufio"
	"net/http"
	"os"
)

const (
	NumPostFields   = 5
	PostSizeInLines = NumPostFields + 2
	PageSizeInLines = PostSizeInLines * 5
)

type QueryParameters struct {
	Page int `json:"page"`
}

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

func ScanPosts(file *os.File, startingLine, endingLine int) []byte {
	if startingLine < 1 {
		startingLine = 1
	}
	if endingLine < startingLine {
		endingLine = startingLine + NumPostFields + 1
	}
	currentLine := 0
	scanner := bufio.NewScanner(file)
	textBytes := make([]byte, 0)
	squareB := []byte("[]")
	textBytes = append(textBytes, squareB[0])
	for scanner.Scan() {
		if currentLine < startingLine {
			currentLine++
			continue
		}
		if currentLine > endingLine {
			break
		}
		textBytes = append(textBytes, scanner.Bytes()...)
		currentLine++
	}
	textBytes = textBytes[:len(textBytes)-1]
	textBytes = append(textBytes, squareB[1])
	return textBytes
}

func PaginatePosts(file *os.File, currentPage int) []byte {
	if currentPage == 0 {
		currentPage = 1
	}
	start := (currentPage-1)*PageSizeInLines + 1
	end := currentPage * PageSizeInLines
	return ScanPosts(file, start, end)
}
