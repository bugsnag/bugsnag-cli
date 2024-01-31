package utils

import (
	"os/user"
)

// GetSystemUser - Gets username from the system
func GetSystemUser() string {
	user, _ := user.Current()
	return user.Username
}
