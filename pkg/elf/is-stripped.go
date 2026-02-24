package elf

import (
	"debug/elf"
	"fmt"
)

// IsStripped checks if an ELF binary has been stripped of symbols and debug information.
//
// A stripped ELF file lacks:
//   - .symtab section (static symbol table)
//   - .debug_* sections (DWARF debug information)
//
// Stripped files retain only .dynsym (dynamic symbols), which are required for runtime linking.
//
// This function is useful for determining whether debug symbols can be extracted from a binary.
// If a file is stripped, tools like objcopy will produce minimal or empty output when extracting
// debug sections.
//
// Parameters:
//
//	filepath - the full file path to the ELF binary.
//
// Returns:
//
//	bool  - true if the file is stripped, false otherwise.
//	error - non-nil if the file cannot be opened or parsed.
func IsStripped(filepath string) (bool, error) {
	file, err := elf.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("failed to open ELF file: %w", err)
	}
	defer file.Close()

	// Check for symbol table section (.symtab)
	// Stripped files only have .dynsym (dynamic symbols)
	if section := file.Section(".symtab"); section != nil {
		return false, nil
	}

	// Also check for debug sections
	for _, section := range file.Sections {
		if len(section.Name) > 6 && section.Name[:6] == ".debug" {
			return false, nil
		}
	}

	return true, nil
}
