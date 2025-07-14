package android_testing

import (
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"os"
	"path/filepath"
	"testing"
)

// TestObjcopy tests that Objcopy correctly extracts architecture and calls objcopy.
func TestObjcopy(t *testing.T) {
	// Create temporary directory for input/output
	tmpDir, err := os.MkdirTemp("", "objcopy-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup mock .so file path like lib/x86/libfoo.so
	libDir := filepath.Join(tmpDir, "lib", "x86")
	if err := os.MkdirAll(libDir, 0755); err != nil {
		t.Fatalf("Failed to create lib dir: %v", err)
	}
	inputFile := filepath.Join(libDir, "libfoo.so")
	if err := os.WriteFile(inputFile, []byte("dummy ELF content"), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	// Create fake objcopy executable
	fakeObjcopy := filepath.Join(tmpDir, "fake-objcopy")
	script := `#!/bin/sh
cp "$3" "$4"
`
	if err := os.WriteFile(fakeObjcopy, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to write fake objcopy: %v", err)
	}

	// Prepend the temp dir to PATH so exec.LookPath finds our fake objcopy
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+origPath)

	// Run Objcopy
	outputPath := filepath.Join(tmpDir, "out")
	outputFile, err := android.Objcopy("fake-objcopy", inputFile, outputPath)
	if err != nil {
		t.Fatalf("Objcopy failed: %v", err)
	}

	// Check that the output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file not created: %s", outputFile)
	}

	// Check expected output path
	expectedOutput := filepath.Join(outputPath, "x86", "libfoo.so.sym")
	if outputFile != expectedOutput {
		t.Errorf("Unexpected output path. Got: %s, Want: %s", outputFile, expectedOutput)
	}
}
