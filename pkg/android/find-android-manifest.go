package android

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// capitalizeFirstLetter returns the input string with the first letter capitalized.
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

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
//   - logger: logger used to emit debug output when manifest is not found.
//
// Returns:
//   - string: The resolved manifest path if found, otherwise an empty string.
func FindAndroidManifest(appBuildPath string, variant string, logger log.Logger) string {
	var err error

	mergedManifestPath := filepath.Join(appBuildPath, "intermediates", "merged_manifests")

	if variant == "" {
		variant, err = GetVariantDirectory(mergedManifestPath)
	}

	if err != nil {
		logger.Info("No AndroidManifest.xml located: a single variant directory couldn't be found")
		return ""
	}

	primaryAppManifestPath := filepath.Join(mergedManifestPath, variant, "AndroidManifest.xml")

	// Fix: Capitalize first letter of variant for fallback path
	capitalizedVariant := capitalizeFirstLetter(variant)
	fallbackAppManifestPath := filepath.Join(mergedManifestPath, variant, "process"+capitalizedVariant+"Manifest", "AndroidManifest.xml")

	paths := []string{primaryAppManifestPath, fallbackAppManifestPath}

	for _, path := range paths {
		if utils.FileExists(path) {
			logger.Debug(fmt.Sprintf("AndroidManifest.xml located at: %s", path))
			return path
		} else {
			logger.Debug(fmt.Sprintf("AndroidManifest.xml not found at: %s", path))
		}
	}

	logger.Info(fmt.Sprintf("No AndroidManifest.xml located for variant %s", variant))
	return ""
}
