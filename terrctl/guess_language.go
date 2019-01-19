package main

import (
	"errors"
	"path/filepath"
	"strings"
)

// GuessLanguage - Guess the programming language according to file extensions
func GuessLanguage(files []WalkedFile) (string, error) {
	for _, file := range files {
		switch ext := strings.ToLower(filepath.Ext(file.Path)); ext {
		case ".c":
			return "c", nil
		case ".rs":
			return "rust", nil
		case ".ts":
			return "assemblyscript", nil
		case ".wasm":
			return "wasm", nil
		}
	}
	return "", errors.New("Unable to detect the programming language")
}
