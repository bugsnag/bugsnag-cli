package elf

import (
	"debug/elf"
	"encoding/hex"
	"fmt"
)

// GetBuildId extracts the GNU Build ID from the given ELF binary.
//
// The Build ID is a unique identifier embedded in ELF files, commonly used
// to associate binaries with debug symbols. This function looks for a NOTE
// segment containing a note with:
//   - Name: "GNU"
//   - Type: 3 (NT_GNU_BUILD_ID)
//
// If found, the function returns the build ID as a hex-encoded string.
// If not found or if the file cannot be read/parsing fails, an error is returned.
//
// Parameters:
//
//	path - the full file path to the ELF binary.
//
// Returns:
//
//	buildID - the hex-encoded build ID string (e.g., "d41d8cd98f00b204e9800998ecf8427e").
//	error   - non-nil if an error occurred or the build ID was not found.
func GetBuildId(filepath string) (string, error) {
	file, err := elf.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open ELF file: %w", err)
	}
	defer file.Close()

	for _, prog := range file.Progs {
		if prog.Type == elf.PT_NOTE {
			data := make([]byte, prog.Filesz)
			_, err := prog.ReadAt(data, 0)
			if err != nil {
				return "", fmt.Errorf("failed to read PT_NOTE segment: %w", err)
			}

			offset := 0
			for offset+12 <= len(data) {
				namesz := int(uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16 | uint32(data[offset+3])<<24)
				descsz := int(uint32(data[offset+4]) | uint32(data[offset+5])<<8 | uint32(data[offset+6])<<16 | uint32(data[offset+7])<<24)
				noteType := uint32(data[offset+8]) | uint32(data[offset+9])<<8 | uint32(data[offset+10])<<16 | uint32(data[offset+11])<<24

				offset += 12
				if offset+namesz > len(data) {
					break
				}
				name := data[offset : offset+namesz]
				offset += ((namesz + 3) & ^3)

				if offset+descsz > len(data) {
					break
				}
				desc := data[offset : offset+descsz]
				offset += ((descsz + 3) & ^3)

				if string(name[:len(name)-1]) == "GNU" && noteType == 3 {
					return hex.EncodeToString(desc), nil
				}
			}
		}
	}

	return "", fmt.Errorf("build ID not found")
}
