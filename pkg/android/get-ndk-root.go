package android

import (
	"fmt"
	"os"
)

// GetAndroidNDKRoot - Returns a valid Android NDK root path
func GetAndroidNDKRoot(path string) (string, error) {

	if path == "" {

		//envValue, envPresent := os.LookupEnv("ANDROID_NDK_ROOT")
		envValue := os.Getenv("ANDROID_NDK_ROOT")

		//if envPresent {

		path = envValue

		//} else {
		//	return "", fmt.Errorf("environment variable 'ANDROID_NDK_ROOT' not defined")
		//}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, fmt.Errorf(path + " does not exist on the system")
	}

	return path, nil
}
