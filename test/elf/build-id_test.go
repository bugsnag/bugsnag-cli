package elf_test

import (
	"github.com/bugsnag/bugsnag-cli/pkg/elf"
	"os"
	"path/filepath"
	"testing"
)

func TestGetBuildId(t *testing.T) {
	// Provide the path to a known ELF binary that has a build ID.
	// You might use something like /bin/ls or a fixture in your testdata directory.
	// For a more controlled test, keep a small compiled ELF binary in ./testdata/

	binPath := filepath.Join("..", "testdata", "unity", "libil2cpp.so")

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Fatalf("Test ELF binary does not exist at path: %s", binPath)
	}

	buildID, err := elf.GetBuildId(binPath)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// You can check against a known build ID if desired
	expected := "c3a9e89c5749add647e4afa0acd484883c47391a"
	if buildID != expected {
		t.Errorf("Expected build ID %s, got %s", expected, buildID)
	}
}
