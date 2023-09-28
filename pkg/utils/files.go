package utils

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
func IsDir(path string) bool {
	pathInfo, err := os.Stat(path)

	return err == nil && pathInfo.IsDir()
}

// BuildFileList - Builds a list of files from a given path(s)
func BuildFileList(paths []string) ([]string, error) {
	var fileList []string

	for _, path := range paths {
		if IsDir(path) {
			log.Info("Searching " + path + " for files")
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

// BuildFolderList - Builds a list of folders from a given path(s)
func BuildFolderList(paths []string) ([]string, error) {
	var folderList []string

	for _, folder := range paths {
		if IsDir(folder) {
			err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
				if err == nil && info.IsDir() {
					if folder != path {
						folderList = append(folderList, filepath.Base(path))
					}
				}
				return nil
			})

			if err != nil {
				return folderList, err
			}
		}
	}
	return folderList, nil
}

// FileExists - Checks if a given file exists on the system
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// FindLatestFileWithSuffix - Finds the latest file with a given suffix
func FindLatestFileWithSuffix(directory string, targetSuffix string) (string, error) {
	var newestFile string
	var newestModTime time.Time

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, targetSuffix) {
			// Check to see if the file that we have found is newer than the previous file
			if info.ModTime().After(newestModTime) {
				newestModTime = info.ModTime()
				newestFile = path
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if newestFile == "" {
		return "", fmt.Errorf("Unable to find " + targetSuffix + " files in " + directory)
	}

	return newestFile, err
}

func ExtractFile(file string, slug string) (string, error) {
	tempDir, err := os.MkdirTemp("", "bugsnag-cli-"+slug+"-unpacking-*")

	if err != nil {
		return "", fmt.Errorf("error creating temporary working directory " + err.Error())
	}

	log.Info("Extracting " + filepath.Base(file) + " to " + tempDir)

	err = Unzip(file, tempDir)

	if err != nil {
		return "", err
	}

	return tempDir, nil
}
