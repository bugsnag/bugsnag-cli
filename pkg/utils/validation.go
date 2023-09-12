package utils

import (
	"os"
	"strings"
)

type UploadPaths []string

// Validate that the path(s) exist
func (p UploadPaths) Validate() error {
	for _, path := range p {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

// XorString checks if one string is not empty and returns a second string if it is
func XorString(string1 string, string2 string) string {
	if string1 != "" {
		return string1
	}
	return string2
}

func ContainsString(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}
