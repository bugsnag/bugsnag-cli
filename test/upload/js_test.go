package upload_testing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
)

func TestPopulateSourceMap(t *testing.T) {
	t.Log("Testing populating source map")
	logger := log.NewLoggerWrapper("debug")

	sourceMapPath := "../testdata/js-nosources/dist/main.js.map"
	results, err := upload.ReadSourceMap(sourceMapPath, logger)
	if err != nil {
		t.Error(err)
	}

	modified := upload.AddSources(results, sourceMapPath, logger)
	if !modified {
		t.Error("Source map should have been modified")
	}

	if len(results["sources"].([]interface{})) != 3 {
		t.Error("Sources is not 3 long")
	}
	contents := results["sourcesContent"].([]*string)
	if len(contents) != 3 {
		t.Error("SourcesContent is not 3 long")
	}
	if !strings.Contains(*contents[2], "const element = document.createElement('div');") {
		t.Error("contents 2 should be populated")
	}
}

func TestProcessJs_IgnoresNodeModulesSourceMaps(t *testing.T) {
	logger := log.NewLoggerWrapper("debug")

	tmpDir := t.TempDir()

	projectRoot := filepath.Join(tmpDir, "project")
	distDir := filepath.Join(projectRoot, "dist")
	nodeModulesDir := filepath.Join(projectRoot, "node_modules", "dep")

	if err := os.MkdirAll(distDir, 0o755); err != nil {
		t.Fatalf("failed to create dist dir: %v", err)
	}
	if err := os.MkdirAll(nodeModulesDir, 0o755); err != nil {
		t.Fatalf("failed to create node_modules dir: %v", err)
	}

	// Sourcemap that should be kept
	validMap := filepath.Join(distDir, "app.js.map")
	if err := os.WriteFile(validMap, []byte(`{"version":3,"sources":[],"names":[],"mappings":""}`), 0o644); err != nil {
		t.Fatalf("failed to write valid sourcemap: %v", err)
	}

	// Sourcemap inside node_modules that should be ignored
	ignoredMap := filepath.Join(nodeModulesDir, "dep.js.map")
	if err := os.WriteFile(ignoredMap, []byte(`{"version":3,"sources":[],"names":[],"mappings":""}`), 0o644); err != nil {
		t.Fatalf("failed to write ignored sourcemap: %v", err)
	}

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			Js: options.Js{
				Path:        []string{distDir},
				ProjectRoot: projectRoot,
				BaseUrl:     "https://example.com/",
			},
		},
	}

	err := upload.ProcessJs(opts, logger)

	// We expect ProcessJs to progress past sourcemap discovery.
	// Any error here should NOT be due to missing sourcemaps.
	if err != nil && strings.Contains(err.Error(), "could not find a source map") {
		t.Fatalf("node_modules sourcemaps were not ignored; no valid sourcemaps detected")
	}
}
