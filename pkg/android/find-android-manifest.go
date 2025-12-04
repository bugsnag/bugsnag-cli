package android

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

// FindAndroidManifest locates and returns the first existing AndroidManifest.xml path.
//
// This function iterates through a list of candidate file paths and returns the
// first one that exists on disk. This is useful when Android build outputs may
// place the manifest in multiple variant-specific directories.
//
// Parameters:
//   - paths: A slice of possible manifest file locations.
//
// Returns:
//   - string: The first existing manifest path, or an empty string if none are found.
func FindAndroidManifest(paths []string) string {
	for _, path := range paths {
		if utils.FileExists(path) {
			return path
		}
	}

	return ""
}
