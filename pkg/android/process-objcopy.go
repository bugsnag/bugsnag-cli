package android

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Objcopy extracts debug symbols from a native shared object (.so) file using `objcopy`.
//
// It infers the architecture from the file path (e.g. "lib/x86/libfoo.so" => "x86"),
// creates an architecture-specific output directory under the given output path,
// and writes the extracted debug symbols to a `.so.sym` file using the `--only-keep-debug` option.
//
// Parameters:
//   - objcopyPath: name or path to the objcopy binary (e.g. "llvm-objcopy").
//   - file: full path to the input .so file.
//   - outputPath: root directory where extracted symbol files should be placed.
//
// Returns:
//   - string: path to the generated `.so.sym` file.
//   - error: non-nil if objcopy fails or if the architecture cannot be inferred.
func Objcopy(objcopyPath string, file string, outputPath string) (string, error) {
	objcopyLocation, err := exec.LookPath(objcopyPath)
	if err != nil {
		return "", err
	}

	// Extract the architecture from the path (e.g. "lib/x86/libpicoapp.so" => "x86")
	var arch string
	parts := strings.Split(file, string(filepath.Separator))
	for i, part := range parts {
		if part == "lib" && i+1 < len(parts) {
			arch = parts[i+1]
			break
		}
	}
	if arch == "" {
		return "", err // Or return a clearer error
	}

	// Create arch-specific output directory
	archDir := filepath.Join(outputPath, arch)
	if err := os.MkdirAll(archDir, os.ModePerm); err != nil {
		return "", err
	}

	// Create the output filename with `.so.sym` extension
	outputFile := filepath.Join(archDir, strings.ReplaceAll(filepath.Base(file), filepath.Ext(file), ".so.sym"))

	_, err = exec.Command(objcopyLocation, "--compress-debug-sections=zlib", "--only-keep-debug", file, outputFile).CombinedOutput()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
