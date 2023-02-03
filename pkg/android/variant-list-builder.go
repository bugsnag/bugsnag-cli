package android

import "io/ioutil"

// BuildVariantsList - Returns a list of variants from a given folder
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
