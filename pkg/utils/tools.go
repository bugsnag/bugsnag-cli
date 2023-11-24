package utils

import (
	"os/exec"
	"strings"
)

const (
	PLUTIL     = "plutil"
	XCODEBUILD = "xcodebuild"
)

// FindLocationOf returns the path of the executable file associated with the given command.
func FindLocationOf(something string) string {
	cmd := exec.Command("which", something)
	location, _ := cmd.Output()
	return strings.TrimSpace(string(location))
}
