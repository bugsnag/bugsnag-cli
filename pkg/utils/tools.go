package utils

import (
	"os/exec"
	"strings"
)

const (
	PLUTIL     = "plutil"
	XCODEBUILD = "xcodebuild"
)

// LocationOf returns the path of the executable file associated with the given command.
func LocationOf(something string) string {
	cmd := exec.Command("which", something)
	location, _ := cmd.Output()
	return strings.TrimSpace(string(location))
}
