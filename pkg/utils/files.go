package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func walk(s string, d fs.DirEntry, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if ! d.IsDir() {
		println(s)
		return s, nil
	}
	return "", nil
}

//FindFilesInDir - Finds files in a given directory
func FindFilesInDir(directory string) ([]string, error) {
	files, err := FilePathWalkDir(directory)
	if err != nil {
		return nil, err
	}
	return files, nil
}

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

	if err != nil {
		return false
	}

	return pathInfo.IsDir()
}