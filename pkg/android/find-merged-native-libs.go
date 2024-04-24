package android

import (
	"fmt"
	"path/filepath"
	
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func findNativeLibPath(arr []string, path string) (string, error) {

	// Look for NDK symbol files if the upload command was run from somewhere within the project root but outside the merged_native_libs directory
	// based on the expected file path of app/build/intermediates/merged_native_libs (or android/app/build/intermediates/merged_native_libs for RN)
	iterations := len(arr)
	for i := 1; i <= iterations; i++ {
		path_ending := filepath.Join(arr...)
		mergeNativeLibPath := filepath.Join(path, path_ending)
		if utils.FileExists(mergeNativeLibPath){
			return mergeNativeLibPath, nil
		}
		arr = arr[1:]
	}

	// Look for NDK symbol files if the upload command was run from somewhere within the merged_native_libs directory
	mergeNativeLibPath := path
	for i := 1; i <= 6; i++ {
		if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
			return mergeNativeLibPath, nil
		}
		mergeNativeLibPath = filepath.Dir(mergeNativeLibPath)
	}

	// If the command was run on the file itself, but a merged_native_libs directory wasn't found in the path to the file
	// set the mergeNativeLibPath based off the file location e.g. merged_native_libs/<variant/out/lib/<arch>/
	if !utils.IsDir(path) {
		mergeNativeLibPath = filepath.Join(path, "..", "..", "..", "..", "..")
		return mergeNativeLibPath, nil
	}

	return "", fmt.Errorf("unable to find the merged_native_libs in " + path)
}