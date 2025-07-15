package utils

import (
	"crypto/md5"
	"fmt"
)

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
