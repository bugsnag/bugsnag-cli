package elf

import (
	"debug/elf"
	"fmt"
)

// GetArch returns the architecture string from the given ELF binary.
//
// This function opens the ELF file and reads its machine architecture field,
// which indicates the target CPU architecture (e.g., "EM_AARCH64", "EM_386").
//
// Parameters:
//
//	path - the full file path to the ELF binary.
//
// Returns:
//
//	arch  - the string representation of the ELF machine architecture.
//	error - non-nil if the file could not be opened or parsed.
func GetArch(filepath string) (string, error) {
	var arch string
	file, err := elf.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open ELF file: %w", err)
	}
	defer file.Close()

	switch file.Machine {
	case elf.EM_AARCH64:
		arch = "arm64"
	case elf.EM_386:
		arch = "x86"
	case elf.EM_X86_64:
		arch = "x86_64"
	case elf.EM_ARM:
		arch = "armv7"
	}

	if arch == "" {
		return file.Machine.String(), fmt.Errorf("unable to find arch type")
	}

	return arch, nil
}
