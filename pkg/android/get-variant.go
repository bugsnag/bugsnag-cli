package android

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// BuildVariantsList - Returns a list of variants from a given path
func BuildVariantsList(path string) ([]string, error) {
	var variants []string

	fileInfo, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		variants = append(variants, file.Name())
	}

	return variants, nil
}

// GetVariantPath - Builds and checks a path with a variant and a file
func GetVariantPath(path string, variant string, file string) (string, string, error) {
	if variant == "" {
		variants, err := BuildVariantsList(path)

		if err != nil {
			return "", "", fmt.Errorf("unable to build list of variants from " + path + " : " + err.Error())
		}

		if len(variants) > 1 {
			return "", "", fmt.Errorf("more than one variant")
		} else {
			variant = variants[0]
		}
	}

	fullPath := filepath.Join(path, variant, file)

	if utils.FileExists(fullPath) {
		return fullPath, variant, nil
	}

	return fullPath, variant, fmt.Errorf(fullPath + " does not exist on the system")
}
