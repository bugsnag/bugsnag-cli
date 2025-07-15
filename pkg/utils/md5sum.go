package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// GetFileMD5 calculates the MD5 checksum of the file at the specified path.
//
// Parameters:
//   - path: File system path to the file.
//
// Returns:
//   - A lowercase hexadecimal string representing the MD5 hash of the file contents.
//   - An error if the file could not be read or the hash could not be computed.
func GetFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to compute MD5 for %q: %w", path, err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetStringMD5 returns the MD5 checksum of the input string.
//
// Parameters:
//   - input: A string to hash.
//
// Returns:
//   - A lowercase hexadecimal string representing the MD5 hash of the input.
func GetStringMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", hash)
}
