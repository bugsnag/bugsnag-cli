package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// ValidatePath - Checks path(s) provided is a valid
func ValidatePath(paths []string) bool {
	for _,path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true
		} else if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return false
}

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

	if pathInfo.IsDir() {
		return true
	}

	return false
}