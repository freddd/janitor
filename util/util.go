package util

import (
	"path/filepath"
	"os"
	"strings"
)

func FindAllFiles(dir string, ignoreDirs []string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		ignore := false
		for _, dir := range ignoreDirs {
			if strings.Contains(path, dir) {
				ignore = true
			}
		}

		if !f.Mode().IsDir() && !ignore {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func CurrentDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return pwd, nil
}