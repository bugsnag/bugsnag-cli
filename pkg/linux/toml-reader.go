package linux

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Manifest Structs must mirror the manifest.toml structure
type Manifest struct {
	Package PackageInfo `toml:"package"`
}

type PackageInfo struct {
	Title   string `toml:"title"`
	Version string `toml:"version"`
	ID      string `toml:"id"`
}

// ReadManifest reads manifest.toml and returns version and id
func ReadManifest(path string) (string, string, error) {
	var manifest Manifest

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read manifest: %w", err)
	}

	// Decode TOML
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return "", "", fmt.Errorf("failed to parse manifest: %w", err)
	}

	return manifest.Package.Version, manifest.Package.ID, nil
}
