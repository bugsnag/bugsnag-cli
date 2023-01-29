package android

import (
	"path/filepath"
	"strconv"
	"strings"
)

// GetNdkVersion - Returns the major NDK version
func GetNdkVersion(path string) (int, error) {

	ndkVersion := strings.Split(filepath.Base(path), ".")

	ndkIntVersion, err := strconv.Atoi(ndkVersion[0])

	if err != nil {
		return 0, err
	}

	return ndkIntVersion, nil
}
