package utils

import (
	"os/user"
)

// GetSystemUser - Gets username from the system
func GetSystemUser() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.Username, nil
}
