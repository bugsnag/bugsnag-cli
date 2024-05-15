package android

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// BuildObjcopyPath - Builds the path to the Objcopy binary within the NDK root path
func BuildObjcopyPath(path string) (string, error) {

	ndkVersion, err := GetNdkVersion(path)

	if err != nil {
		return "", fmt.Errorf("unable to determine ndk version from %s", path)
	}

	if ndkVersion < 24 {
		return "", fmt.Errorf("unsupported NDK version. Please upgrade to r24 or higher")
	} else {
		directoryPattern := filepath.Join(path, "/toolchains/llvm/prebuilt/*/bin")

		directoryMatches, err := filepath.Glob(directoryPattern)

		if err != nil {
			return "", err
		}

		if directoryMatches == nil {
			return "", fmt.Errorf("Unable to find objcopy within: %s", path)
		}

		if runtime.GOOS == "windows" {
			return filepath.Join(directoryMatches[0], "llvm-objcopy.exe"), nil
		}

		return filepath.Join(directoryMatches[0], "llvm-objcopy"), nil
	}
}
