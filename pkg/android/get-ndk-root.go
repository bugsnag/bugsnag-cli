package android

import (
	"fmt"
	"os"
)

// GetAndroidNDKRoot - Returns a valid Android NDK root path
func GetAndroidNDKRoot(path string) (string, error) {

	if path == "" {
		envValue, envPresent := os.LookupEnv("ANDROID_NDK_ROOT")

		if envPresent {

		} else {
			return "", fmt.Errorf("environment variable 'ANDROID_NDK_ROOT' not defined")
		}

		path = envValue
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, fmt.Errorf("%s does not exist on the system", path)
	}

	return path, nil
}
