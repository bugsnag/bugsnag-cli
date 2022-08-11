package utils

import (
	"errors"
	"os"
)

// ValidatePath Checks if a provided string is a valid path or file
func ValidatePath(paths []string) (bool, string){
	for _,path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true, path
		} else if errors.Is(err, os.ErrNotExist) {
			return false, path
		}
	}
	return false, ""
}


