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

func GetVariant(path string) (string, error) {
	var variants []string

	fileInfo, err := ioutil.ReadDir(path)

	if err != nil {
		return "", err
	}

	for _, file := range fileInfo {
		variants = append(variants, file.Name())
	}

	if len(variants) > 1 {
		return "", fmt.Errorf("more than one variant found. Please specify using `--variant` ")
	} else if len(variants) < 1 {
		return "", fmt.Errorf("no variants found. Please specify using `--variant`")
	}

	variant := variants[0]

	if !utils.FileExists(filepath.Join(path, variant)) {
		return "", fmt.Errorf("variant path " + filepath.Join(path, variant) + " doesn't exist on the system")
	}

	return variant, nil
}
