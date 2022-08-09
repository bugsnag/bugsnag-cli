package validation

import (
	"errors"
	"os"
)

// ValidatePath Checks if a provided string is a valid path or file
func ValidatePath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}