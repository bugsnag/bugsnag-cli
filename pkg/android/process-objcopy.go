package android

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Objcopy - Processes files using objcopy
func Objcopy(objcopyPath string, file string) (string, error) {

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "")

	if err != nil {
		return "", err
	}

	objcopyLocation, err := exec.LookPath(objcopyPath)

	if err != nil {
		return "", err
	}

	outputFile := filepath.Join(tempDir, filepath.Base(file))
	outputFile = strings.ReplaceAll(outputFile, filepath.Ext(outputFile), ".so.sym")

	cmd := exec.Command(objcopyLocation, "--compress-debug-sections=zlib", "--only-keep-debug", file, outputFile)

	_, err = cmd.CombinedOutput()

	return outputFile, nil
}
