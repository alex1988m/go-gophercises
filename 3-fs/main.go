package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

func main() {
	err := processFiles(".", filenameFilter, fileHandle)
	if err != nil {
		logrus.Fatal(err)
	}
}

func fileHandle(path string) error {
	dir := filepath.Dir(path)
	oldName := filepath.Base(path)

	// Extract xxx and yyy from the filename
	parts := regexp.MustCompile(`^(.+)_(\d+)\.txt$`).FindStringSubmatch(oldName)
	if len(parts) != 3 {
		return fmt.Errorf("unexpected filename format: %s", oldName)
	}

	xxx := parts[1]
	yyy := parts[2]

	// Create the new filename
	newName := fmt.Sprintf("%s (%s of 100).txt", xxx, yyy)
	newPath := filepath.Join(dir, newName)

	// Rename the file
	err := os.Rename(path, newPath)
	if err != nil {
		return fmt.Errorf("error renaming file: %v", err)
	}

	logrus.Infof("Renamed: %s -> %s", oldName, newName)
	return nil
}

func processFiles(root string, filter func(string) bool, handler func(string) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filter(d.Name()) {
			if err := handler(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func filenameFilter(filename string) bool {
	// Regex pattern to match xxx_yyy.txt where yyy are numbers
	pattern := `^[a-zA-Z]+_\d+\.txt$`
	match, _ := regexp.MatchString(pattern, filename)
	return match
}
