package android

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// BuildVariantsList - Returns a list of variants from a given path
func BuildVariantsList(path string) ([]string, error) {
	var variants []string

	fileInfo, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		variants = append(variants, file.Name())
	}

	return variants, nil
}

func GetVariantDirectory(path string) (string, error) {
	var variants []string

	fileInfo, err := os.ReadDir(path)

	if err != nil {
		return "", err
	}

	for _, file := range fileInfo {
		variants = append(variants, file.Name())
	}

	if len(variants) > 1 {
		return "", fmt.Errorf("more than one variant found. Please specify using `--variant`")
	} else if len(variants) < 1 {
		return "", fmt.Errorf("no variants found. Please specify using `--variant`")
	}

	variant := variants[0]

	if !utils.FileExists(filepath.Join(path, variant)) {
		return "", fmt.Errorf("variant path %s doesn't exist on the system", filepath.Join(path, variant))
	}

	return variant, nil
}

func FindVariantDexFiles(mappingFilePath string, variant string) []string {
	buildRoot := filepath.Join(filepath.Dir(mappingFilePath), "..", "..", "..", "intermediates", "dex", variant)

	if utils.IsDir(buildRoot) {
		matches, _ := filepath.Glob(filepath.Join(buildRoot, "*", "classes.dex"))
		return matches
	}

	return []string{}
}
