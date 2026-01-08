package ios

import (
	"encoding/binary"
	"os"
)

// AppleDouble magic number (big-endian)
const appleDoubleMagic = 0x00051607

// IsAppleDoubleMetaData checks whether a given file is an AppleDouble metadata file.
//
// Parameters:
// - path: The file path to examine.
//
// Returns:
// - A boolean indicating whether the file is AppleDouble metadata.
// - An error if the file cannot be opened or read.
func IsAppleDoubleMetaData(path string) (bool, error) {
	f, err := os.Open(path)

	if err != nil {
		return false, err
	}

	defer f.Close()

	var magic uint32
	if err := binary.Read(f, binary.BigEndian, &magic); err != nil {
		return false, err
	}

	return magic == appleDoubleMagic, nil
}
