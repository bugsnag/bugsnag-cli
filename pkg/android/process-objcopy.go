package android

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// Objcopy extracts debug symbols from a native library using `objcopy` and compresses them.
//
// It computes the MD5 hash of the input file to generate a consistent output filename
// and uses `objcopy` to write only debug information into a compressed `.sym` file.
//
// Parameters:
//   - objcopyPath: Path to the `objcopy` binary (from Android NDK).
//   - file: Path to the native library (.so) to process.
//   - outputPath: Directory where the resulting .sym file should be saved.
//
// Returns:
//   - string: Full path to the generated .sym file.
//   - error: Non-nil if the operation fails at any point (e.g., invalid objcopy path, command failure).
func Objcopy(objcopyPath, file, outputPath string) (string, error) {
	objcopyLocation, err := exec.LookPath(objcopyPath)
	if err != nil {
		return "", fmt.Errorf("objcopy binary not found at path %s: %w", objcopyPath, err)
	}

	md5sum, err := utils.GetFileMD5(file)
	if err != nil {
		return "", fmt.Errorf("failed to calculate MD5 of file %s: %w", file, err)
	}

	outputFile := filepath.Join(outputPath, md5sum)

	cmd := exec.Command(objcopyLocation, "--compress-debug-sections=zlib", "--only-keep-debug", file, outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("objcopy failed for file %s: %w\nCommand output: %s", file, err, string(output))
	}

	return outputFile, nil
}
