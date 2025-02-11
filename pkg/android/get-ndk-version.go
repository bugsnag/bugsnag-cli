package android

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetNdkVersion - Returns the major NDK version
func GetNdkVersion(ndkPath string) (int, error) {
	// Construct the path to source.properties
	sourceFile := filepath.Join(ndkPath, "source.properties")

	// Open the file
	file, err := os.Open(sourceFile)
	if err != nil {
		return 0, fmt.Errorf("failed to open source.properties: %v", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Pkg.Revision") {
			// Extract version string
			parts := strings.Split(line, "=")
			parts[1] = strings.TrimSpace(parts[1])
			if len(parts) == 2 {
				versionStr := parts[1]
				// Extract major version (before the first dot)
				versionParts := strings.Split(versionStr, ".")
				if len(versionParts) > 0 {
					majorVersion, err := strconv.Atoi(versionParts[0])
					if err != nil {
						return 0, fmt.Errorf("failed to parse major version: %v", err)
					}
					return majorVersion, nil
				}
			}
		}
	}

	// Check for errors while scanning
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading source.properties: %v", err)
	}

	return 0, fmt.Errorf("NDK version not found in source.properties")
}
