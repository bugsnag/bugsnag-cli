package utils

import (
	"fmt"
	"os"
	"strings"
)

type Paths []string
type Path string
type LogLevels string

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

// ContainsString Check if a string is in a slice
func ContainsString(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}

// Validate that the log level is valid
func (l LogLevels) Validate() error {
	switch strings.ToLower(string(l)) {
	case "debug", "info", "warn", "fatal":
		return nil
	default:
		return fmt.Errorf("invalid log level: %s", l)
	}
}
