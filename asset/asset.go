package asset

import (
	"bufio"
	"embed"
	"fmt"
	"path"
	"strings"
)

//go:embed word_list/*
var files embed.FS

const WordListKey = "word_list"

var fileMap = map[string]map[string]string{
	WordListKey: {
		"EN": "en.txt",
	},
}

func keyToFile(key, fileType string) (string, bool) {
	file, ok := fileMap[fileType][strings.ToUpper(key)]

	return file, ok
}

func getWordListFilePath(key string) (string, error) {
	fileName, ok := keyToFile(key, WordListKey)
	if !ok {
		return "", fmt.Errorf("invalid %s value (%s)", WordListKey, key)
	}

	return path.Join(WordListKey, fileName), nil
}

// GetFilteredWordList reads a word list from an embedded file identified by the
// given key, and filters the words based on the specified length. It returns a
// slice of strings that meet the length criteria. If the file cannot be opened
// or read, or if an error occurs during scanning, an error is returned.
func GetFilteredWordList(key string, minMaxLen int) ([]string, error) {
	path, err := getWordListFilePath(key)
	if err != nil {
		return nil, err
	}

	file, err := files.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wl []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == minMaxLen {
			wl = append(wl, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wl, nil
}
