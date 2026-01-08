package ios

import (
	"encoding/binary"
	"os"
)

// AppleDouble magic number (big-endian)
const appleDoubleMagic = 0x00051607

// IsAppleDoubleMetaData reports whether the given file is an AppleDouble metadata file.
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
