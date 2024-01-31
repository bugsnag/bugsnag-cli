package utils

import (
	"os"
	"strings"
)

type Paths []string
type Path string

// Validate that the path(s) exist
func (p Paths) Validate() error {
	for _, path := range p {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

// Validate that the path exist
func (p Path) Validate() error {
	if _, err := os.Stat(string(p)); err != nil {
		return err
	}
	return nil
}

func ContainsString(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}
