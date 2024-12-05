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

// GetXcodeArchiveLocation retrieves the custom archive location for Xcode archives from user preferences.
// It first attempts to read the custom location from the system defaults. If the location is not set,
// it returns the default Xcode archive location within the user's home directory.
//
// Parameters:
// - None
//
// Returns:
// - A string containing the path to the Xcode archive location.
// - An error if the location cannot be determined or if an error occurs during command execution.
func GetXcodeArchiveLocation() (string, error) {
	// Command to read the custom archive location from Xcode preferences
	cmd := exec.Command("defaults", "read", "com.apple.dt.Xcode", "IDECustomDistributionArchivesLocation")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command and capture the output
	err := cmd.Run()
	if err != nil {
		// If the command fails, check stderr for additional details
		if strings.Contains(stderr.String(), "does not exist") {
			// If the key is not set, return the default location for Xcode archives
			usr, err := user.Current()
			if err != nil {
				return "", fmt.Errorf("unable to get current user: %w", err)
			}
			// Return default archive location under user's home directory
			return filepath.Join(usr.HomeDir, "Library", "Developer", "Xcode", "Archives"), nil
		}
		// Return error with additional details from stderr if the command fails
		return "", fmt.Errorf("error running defaults command: %w, stderr: %s", err, stderr.String())
	}

	// Trim any trailing newline or spaces from the output
	customPath := strings.TrimSpace(out.String())
	if customPath == "" {
		// Return error if the command succeeded but no path was returned
		return "", fmt.Errorf("command succeeded but returned an empty path")
	}

	return customPath, nil
}

// GetLatestXcodeArchive finds the latest .xcarchive file for a given scheme in the specified directory.
// It walks through the directory, looking for .xcarchive files that match the scheme. The most recently modified
// archive file is returned.
//
// Parameters:
// - path: Directory path to search for .xcarchive files.
// - scheme: The scheme used as a prefix to filter relevant archive files.
//
// Returns:
// - A string containing the path to the most recent .xcarchive file matching the scheme.
// - An error if no matching archive is found or if an error occurs during the search.
func GetLatestXcodeArchive(path, scheme string) (string, error) {
	// Attempt to find the latest .xcarchive for the given scheme in the folder
	path, err := findLatestXCArchive(path, scheme)
	if err != nil {
		return "", err
	} else if path == "" {
		// If no archive is found, return an error indicating the folder doesn't contain any matching archives
		return "", fmt.Errorf("No xcarchive found in %s", path)
	}
	// Return the path of the latest archive found
	return path, nil
}

// findLatestXCArchive walks through the folder and looks for the latest .xcarchive file that matches the scheme.
// It compares modification times and returns the path of the most recently modified archive.
//
// Parameters:
// - folderPath: Directory path to search for .xcarchive files.
// - scheme: The scheme used as a prefix to filter relevant archive files.
//
// Returns:
// - A string containing the path to the most recent .xcarchive file matching the scheme.
// - An error if an issue occurs during the directory walk or if no matching archive is found.
func findLatestXCArchive(folderPath, scheme string) (string, error) {
	var latestFile string
	var latestModTime time.Time

	// Walk through all files and directories in the folder
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If an error occurs while walking, return it
			return err
		}

		// Check if the entry is a directory, ends with ".xcarchive", and starts with the given scheme
		if info.IsDir() && strings.HasSuffix(info.Name(), ".xcarchive") && strings.HasPrefix(info.Name(), scheme) {
			// If this archive is newer than the previous one, update the latestFile and latestModTime
			if info.ModTime().After(latestModTime) {
				latestFile = path
				latestModTime = info.ModTime()
			}
		}
		return nil
	})

	// If an error occurred during the walk, return it
	if err != nil {
		return "", err
	}

	// If no matching .xcarchive file was found, return an error
	if latestFile == "" {
		return "", fmt.Errorf("no .xcarchive files found for scheme: %s", scheme)
	}

	// Return the path of the latest archive
	return latestFile, nil
}
