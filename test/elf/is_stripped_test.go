package elf_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/elf"
	"github.com/stretchr/testify/assert"
)

func TestIsStripped_StrippedFile(t *testing.T) {
	// This file is confirmed to be stripped (no .symtab, no .debug_* sections)
	strippedFile := "../../features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/armeabi-v7a/libbugsnag-ndk.so"

	isStripped, err := elf.IsStripped(strippedFile)
	assert.NoError(t, err, "Should be able to check if file is stripped")
	assert.True(t, isStripped, "File should be detected as stripped")
}

func TestIsStripped_UnstrippedFile(t *testing.T) {
	// This file is confirmed to have debug symbols (not stripped)
	unstrippedFile := "../../platforms-examples/Unity/Library/Bee/artifacts/Android/libunity/armeabi-v7a/unstripped/libunity.so"

	isStripped, err := elf.IsStripped(unstrippedFile)
	assert.NoError(t, err, "Should be able to check if file is stripped")
	assert.False(t, isStripped, "File should be detected as not stripped")
}

func TestIsStripped_NonExistentFile(t *testing.T) {
	// Test with a file that doesn't exist
	nonExistentFile := "/path/to/nonexistent/file.so"

	_, err := elf.IsStripped(nonExistentFile)
	assert.Error(t, err, "Should return error for non-existent file")
}

func TestIsStripped_InvalidElfFile(t *testing.T) {
	// Test with a non-ELF file
	invalidFile := "../../CHANGELOG.md"

	_, err := elf.IsStripped(invalidFile)
	assert.Error(t, err, "Should return error for non-ELF file")
}
