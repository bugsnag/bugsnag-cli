package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	PLUTIL     = "plutil"
	XCODEBUILD = "xcodebuild"
	DWARFDUMP  = "dwarfdump"
)

// FilePathWalkDir recursively finds all files within a given directory.
//
// Parameters:
// - root (string): The root directory to start searching from.
//
// Returns:
// - []string: A slice of file paths found within the directory.
// - error: Any error encountered during the walk process.
func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// IsDir checks if the provided path is a directory.
//
// Parameters:
// - path (string): The path to check.
//
// Returns:
// - bool: True if the path is a directory; false otherwise.
func IsDir(path string) bool {
	pathInfo, err := os.Stat(path)
	return err == nil && pathInfo.IsDir()
}

// BuildFileList compiles a list of files from the provided paths.
//
// Parameters:
// - paths ([]string): A slice of paths to process.
//
// Returns:
// - []string: A slice containing file paths from directories and standalone files.
// - error: Any error encountered during processing.
func BuildFileList(paths []string) ([]string, error) {
	var fileList []string

	for _, path := range paths {
		if IsDir(path) {
			files, err := FilePathWalkDir(path)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, files...)
		} else {
			fileList = append(fileList, path)
		}
	}

	return fileList, nil
}

// BuildDirectoryList compiles a list of directories from the provided paths.
//
// Parameters:
// - paths ([]string): A slice of paths to process.
//
// Returns:
// - []string: A slice containing the base names of subdirectories.
// - error: Any error encountered during processing.
func BuildDirectoryList(paths []string) ([]string, error) {
	var directoryList []string

	for _, directory := range paths {
		if IsDir(directory) {
			err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && directory != path {
					directoryList = append(directoryList, filepath.Base(path))
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return directoryList, nil
}

// FileExists checks if a given file exists.
//
// Parameters:
// - path (string): The file path to check.
//
// Returns:
// - bool: True if the file exists; false otherwise.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FindLatestFileWithSuffix searches for the most recently modified file with a given suffix.
//
// Parameters:
// - directory (string): The directory to search in.
// - targetSuffix (string): The suffix to match.
//
// Returns:
// - string: The path to the newest file matching the suffix.
// - error: Any error encountered during the search.
func FindLatestFileWithSuffix(directory, targetSuffix string) (string, error) {
	var newestFile string
	var newestModTime time.Time

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, targetSuffix) && info.ModTime().After(newestModTime) {
			newestModTime = info.ModTime()
			newestFile = path
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	if newestFile == "" {
		return "", fmt.Errorf("unable to find files with suffix '%s' in '%s'", targetSuffix, directory)
	}

	return newestFile, nil
}

// FindLatestFolder searches for the most recently modified folder in a given directory.
//
// Parameters:
// - directory (string): The directory to search in.
//
// Returns:
// - string: The path to the newest folder.
// - error: Any error encountered during the search.
func FindLatestFolder(directory string) (string, error) {
	var newestFolder string
	var newestModTime time.Time

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.ModTime().After(newestModTime) {
			newestModTime = info.ModTime()
			newestFolder = path
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	if newestFolder == "" {
		return "", fmt.Errorf("unable to find folders in '%s'", directory)
	}

	return newestFolder, nil
}

// ExtractFile extracts the contents of a file into a temporary directory.
//
// Parameters:
// - file (string): The file to extract.
// - slug (string): A unique identifier for the temporary directory.
//
// Returns:
// - string: The path to the temporary directory containing the extracted files.
// - error: Any error encountered during extraction.
func ExtractFile(file, slug string) (string, error) {
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("bugsnag-cli-%s-unpacking-*", slug))
	if err != nil {
		return "", fmt.Errorf("error creating temporary working directory: %s", err.Error())
	}

	if err := Unzip(file, tempDir); err != nil {
		return "", err
	}

	return tempDir, nil
}

// FindFolderWithSuffix searches for the first folder with a specified suffix.
//
// Parameters:
// - rootPath (string): The root directory to search in.
// - targetSuffix (string): The suffix to match.
//
// Returns:
// - string: The path to the matching folder.
// - error: Any error encountered during the search.
func FindFolderWithSuffix(rootPath, targetSuffix string) (string, error) {
	var matchingFolder string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasSuffix(info.Name(), targetSuffix) {
			matchingFolder = path
			return filepath.SkipDir
		}
		return nil
	})

	return matchingFolder, err
}

// LocationOf determines the path of the executable associated with a given command.
//
// Parameters:
// - something (string): The command to locate.
//
// Returns:
// - string: The path to the executable or an empty string if not found.
func LocationOf(something string) string {
	cmd := exec.Command("which", something)
	location, _ := cmd.Output()
	return strings.TrimSpace(string(location))
}
