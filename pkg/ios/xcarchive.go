package ios

import (
	"bytes"
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// GetLatestXcodeArchiveForScheme retrieves the latest xcarchive for a given scheme.
// It determines the archive location (either custom or default) and then searches
// for the most recently modified xcarchive matching the scheme.
//
// Parameters:
// - scheme: The scheme used as a prefix to filter relevant archive files.
//
// Returns:
// - A string containing the path to the most recent xcarchive file matching the scheme.
// - An error if the location cannot be determined or if no matching archive is found.
func GetLatestXcodeArchiveForScheme(scheme string) (string, error) {
	// Retrieve the xcarchive location
	archivePath, err := func() (string, error) {
		// Command to read the custom archive location from Xcode preferences
		cmd := exec.Command("defaults", "read", "com.apple.dt.Xcode", "IDECustomDistributionArchivesLocation")

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		// Execute the command and capture the output
		if err := cmd.Run(); err != nil {
			// If the command fails, check stderr for additional details
			if strings.Contains(stderr.String(), "does not exist") {
				// If the key is not set, return the default location for xcarchives
				usr, err := user.Current()
				if err != nil {
					return "", fmt.Errorf("unable to get current user: %w", err)
				}
				// Return default archive location under user's home directory
				archivePath := filepath.Join(usr.HomeDir, "Library", "Developer", "Xcode", "Archives")
				if utils.IsDir(archivePath) {
					return archivePath, nil
				}
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
	}()

	if err != nil {
		return "", fmt.Errorf("failed to determine Xcode archive location: %w", err)
	}

	// Search for the latest xcarchive matching the scheme
	var latestFile string
	var latestModTime time.Time

	err = filepath.Walk(archivePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			// If an error occurs while walking, return it
			return err
		}

		// Check if the entry is a directory, ends with ".xcarchive", and starts with the given scheme
		if info.IsDir() && strings.HasSuffix(info.Name(), ".xcarchive") && strings.HasPrefix(info.Name(), scheme) {
			// If this archive is newer than the previous one, update the latestFile and latestModTime
			if info.ModTime().After(latestModTime) {
				latestFile = filePath
				latestModTime = info.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking the Xcode archive directory: %w", err)
	}

	// If no matching xcarchive file was found, return an error
	if latestFile == "" {
		return "", fmt.Errorf("no Xcode archive files found for scheme: %s", scheme)
	}

	return latestFile, nil
}

// FindXcarchivePath searches for an Xcode archive (.xcarchive) within the provided paths.
// If a directory is given, it attempts to locate the latest archive associated with a scheme.
//
// Parameters:
// - options: CLI options containing upload paths and shared settings.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - The path to the found Xcode archive (.xcarchive), or an empty string if none are found.
// - An error if any issue occurs during the process.
func FindXcarchivePath(options options.CLI, logger log.Logger) (string, error) {
	var (
		xcarchivePath string
		err           error
	)

	// Iterate through the list of provided paths
	for _, path := range options.Upload.XcodeArchive.Path {
		if filepath.Ext(path) == ".xcarchive" {
			// If the path is already an Xcode archive, assign it
			xcarchivePath = path
		} else if utils.IsDir(path) {
			// If the path is a directory, search for an Xcode archive inside it
			logger.Info(fmt.Sprintf("Searching for an Xcode archive in %s", path))

			// Check if the directory contains an Xcode project or workspace
			if IsPathAnXcodeProjectOrWorkspace(path) {
				// Determine the scheme if it is not set
				if options.Upload.XcodeArchive.Shared.Scheme == "" {
					options.Upload.XcodeArchive.Shared.Scheme, err = GetDefaultScheme(path)
					if err != nil {
						return "", fmt.Errorf("error determining default scheme: %w", err)
					}
				}

				// Retrieve the latest Xcode archive for the determined scheme
				xcarchivePath, _ = GetLatestXcodeArchiveForScheme(options.Upload.XcodeArchive.Shared.Scheme)
			}
		}
	}

	return xcarchivePath, nil
}
