package upload_testing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
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

func TestResolveSourceMapPaths_IgnoresNodeModulesAndCssMaps(t *testing.T) {
	logger := log.NewLoggerWrapper("debug")

	// Create a temporary directory structure
	tempDir := t.TempDir()

	// Valid source map
	validMap := filepath.Join(tempDir, "app.js.map")
	if err := os.WriteFile(validMap, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write valid source map: %v", err)
	}

	// CSS source map (should be ignored)
	cssMap := filepath.Join(tempDir, "styles.css.map")
	if err := os.WriteFile(cssMap, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write css source map: %v", err)
	}

	// node_modules source map (should be ignored)
	nodeModulesDir := filepath.Join(tempDir, "node_modules", "lib")
	if err := os.MkdirAll(nodeModulesDir, 0755); err != nil {
		t.Fatalf("failed to create node_modules dir: %v", err)
	}
	nodeModulesMap := filepath.Join(nodeModulesDir, "lib.js.map")
	if err := os.WriteFile(nodeModulesMap, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write node_modules source map: %v", err)
	}

	paths, err := upload.ResolveSourceMapPaths("", tempDir, logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("expected 1 source map, got %d", len(paths))
	}

	if paths[0] != validMap {
		t.Errorf("expected source map %s, got %s", validMap, paths[0])
	}
}
