package ios

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func GetXcodeArchiveLocation() (string, error) {
	// Command to read the custom archive location
	cmd := exec.Command("defaults", "read", "com.apple.dt.Xcode", "IDECustomDistributionArchivesLocation")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		// If the command fails, check stderr for additional details
		if strings.Contains(stderr.String(), "does not exist") {
			// If key is not set, return the default location
			usr, err := user.Current()
			if err != nil {
				return "", fmt.Errorf("unable to get current user: %w", err)
			}
			return filepath.Join(usr.HomeDir, "Library", "Developer", "Xcode", "Archives"), nil
		}
		return "", fmt.Errorf("error running defaults command: %w, stderr: %s", err, stderr.String())
	}

	// Trim any trailing newline or spaces from the output
	customPath := strings.TrimSpace(out.String())
	if customPath == "" {
		return "", fmt.Errorf("command succeeded but returned an empty path")
	}

	return customPath, nil
}

// GetLatestXcodeArchive Finds the latest .xcarchive file in a given directory + YYYY-MM-DD format
func GetLatestXcodeArchive(path, scheme string) (string, error) {
	path, err := findLatestXCArchive(path, scheme)
	if err != nil {
		return "", err
	} else if path == "" {
		return "", fmt.Errorf("No xcarchive found in %s", path)
	}
	return path, nil
}

func findLatestXCArchive(folderPath, scheme string) (string, error) {
	var latestFile string
	var latestModTime time.Time

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error from filepath.Walk
		}

		// Check if the file matches the criteria
		if info.IsDir() && strings.HasSuffix(info.Name(), ".xcarchive") && strings.HasPrefix(info.Name(), scheme) {
			if info.ModTime().After(latestModTime) {
				latestFile = path
				latestModTime = info.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		return "", err // Return any error encountered during walking
	}

	if latestFile == "" {
		return "", fmt.Errorf("no .xcarchive files found for scheme: %s", scheme)
	}

	return latestFile, nil
}
