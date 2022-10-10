package utils

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"os"
	"path/filepath"
)

// FilePathWalkDir - finds files within a given directory
func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// IsDir - Checks if a provided path is a directory or not
func IsDir(path string) bool{
	pathInfo, err := os.Stat(path)

	return err == nil && pathInfo.IsDir()
}

// BuildFileList - Builds a list of files from a given path(s)
func BuildFileList(paths []string) ([]string, error) {
	var fileList []string

	for _, path := range paths {
		if IsDir(path) {
			log.Info("Searching directory for files")
			files, err := FilePathWalkDir(path)
			if err != nil {
				return nil, err
			}
			for _, s := range files {
				fileList = append(fileList, s)
			}
		} else {
			fileList = append(fileList, path)
		}
	}

	return fileList, nil
}
