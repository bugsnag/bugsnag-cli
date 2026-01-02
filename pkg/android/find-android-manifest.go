package android

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// FindAndroidManifest locates and returns the first existing AndroidManifest.xml
// for a given build directory and variant.
//
// This function constructs and checks multiple potential manifest locations that
// vary based on Android build outputs. It returns the first manifest file found
// on disk.
//
// Parameters:
//   - appBuildPath: The root build directory where Android build outputs are located.
//   - variant: The build variant (e.g., "debug", "release") whose manifest paths
//     should be searched.
//
// Returns:
//   - string: The resolved manifest path if found, otherwise an empty string.
//   - error: Non-nil if the manifest cannot be found.
func FindAndroidManifest(appBuildPath string, variant string) (string, error) {
	var err error

	mergedManifestPath := filepath.Join(appBuildPath, "intermediates", "merged_manifests")

	if variant == "" {
		variant, err = GetVariantDirectory(mergedManifestPath)
	}

	if err != nil {
		return "", err
	}

	primaryAppManifestPath := filepath.Join(mergedManifestPath, variant, "AndroidManifest.xml")

	fallbackAppManifestPath := filepath.Join(mergedManifestPath, variant, "process"+variant+"Manifest", "AndroidManifest.xml")

	paths := []string{primaryAppManifestPath, fallbackAppManifestPath}

	for _, path := range paths {
		if utils.FileExists(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("unable to locate AndroidManifest.xml for variant %s", variant)
}
