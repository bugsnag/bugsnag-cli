package macho

import (
	"debug/macho"
	"fmt"
)

// GetArch returns the architecture string from the given thin Mach-O file.
//
// This function opens the Mach-O file and reads the CPU type from its header,
// returning the architecture name in the form expected by dwarfdump's --arch
// flag (e.g., "arm64", "x86_64").
//
// Parameters:
//
//	path - the full file path to the Mach-O file.
//
// Returns:
//
//	arch  - the architecture name of the Mach-O file.
//	error - non-nil if the file could not be opened or parsed, or if its CPU
//	        type is not supported.
func GetArch(path string) (string, error) {
	file, err := macho.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open Mach-O file: %w", err)
	}
	defer file.Close()

	switch file.Cpu {
	case macho.CpuArm64:
		return "arm64", nil
	case macho.CpuAmd64:
		return "x86_64", nil
	}

	return file.Cpu.String(), fmt.Errorf("unable to find arch type")
}
