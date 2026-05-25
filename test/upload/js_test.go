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

	// Create a valid bundle with sourceMappingURL
	validBundle := filepath.Join(tempDir, "app.js")
	bundleContent := "console.log('test');\n//# sourceMappingURL=app.js.map\n"
	if err := os.WriteFile(validBundle, []byte(bundleContent), 0644); err != nil {
		t.Fatalf("failed to write valid bundle: %v", err)
	}

	// Valid source map
	validMap := filepath.Join(tempDir, "app.js.map")
	if err := os.WriteFile(validMap, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write valid source map: %v", err)
	}

	// CSS source map (should be ignored - no .css bundle)
	cssMap := filepath.Join(tempDir, "styles.css.map")
	if err := os.WriteFile(cssMap, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write css source map: %v", err)
	}

	// node_modules bundle and source map (should be ignored)
	nodeModulesDir := filepath.Join(tempDir, "node_modules", "lib")
	if err := os.MkdirAll(nodeModulesDir, 0755); err != nil {
		t.Fatalf("failed to create node_modules dir: %v", err)
	}
	nodeModulesBundle := filepath.Join(nodeModulesDir, "lib.js")
	nodeModulesBundleContent := "console.log('lib');\n//# sourceMappingURL=lib.js.map\n"
	if err := os.WriteFile(nodeModulesBundle, []byte(nodeModulesBundleContent), 0644); err != nil {
		t.Fatalf("failed to write node_modules bundle: %v", err)
	}
	nodeModulesMap := filepath.Join(nodeModulesDir, "lib.js.map")
	if err := os.WriteFile(nodeModulesMap, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write node_modules source map: %v", err)
	}

	bundles, err := upload.ResolveSourceMapPaths("", "", tempDir, logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(bundles) != 1 {
		t.Fatalf("expected 1 source map bundle, got %d", len(bundles))
	}

	if bundles[0].BundlePath != validBundle {
		t.Errorf("expected bundle %s, got %s", validBundle, bundles[0].BundlePath)
	}

	if bundles[0].SourceMapPath != validMap {
		t.Errorf("expected source map %s, got %s", validMap, bundles[0].SourceMapPath)
	}
}

func TestExtractSourceMappingURL(t *testing.T) {
	logger := log.NewLoggerWrapper("debug")
	tempDir := t.TempDir()

	t.Run("Modern syntax with //#", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "modern.js")
		content := "console.log('test');\n//# sourceMappingURL=modern.js.map\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "modern.js.map" {
			t.Errorf("expected 'modern.js.map', got '%s'", url)
		}
	})

	t.Run("Legacy syntax with //@", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "legacy.js")
		content := "console.log('test');\n//@ sourceMappingURL=legacy.js.map\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "legacy.js.map" {
			t.Errorf("expected 'legacy.js.map', got '%s'", url)
		}
	})

	t.Run("No sourceMappingURL", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "no-url.js")
		content := "console.log('test');\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "" {
			t.Errorf("expected empty string, got '%s'", url)
		}
	})

	t.Run("Relative path", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "relative.js")
		content := "console.log('test');\n//# sourceMappingURL=../maps/relative.js.map\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "../maps/relative.js.map" {
			t.Errorf("expected '../maps/relative.js.map', got '%s'", url)
		}
	})

	t.Run("Data URL (inline source map)", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "inline.js")
		content := "console.log('test');\n//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozfQ==\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasPrefix(url, "data:") {
			t.Errorf("expected data URL, got '%s'", url)
		}
	})

	// TC39 ECMA-426 Spec Compliance Tests
	t.Run("Ignores sourceMappingURL in string literal", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "in-string.js")
		content := "let a = \"//# sourceMappingURL=fake.js.map\";\nconsole.log(a);\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should be empty because sourceMappingURL inside a string literal must be ignored
		if url != "" {
			t.Errorf("expected empty string (should ignore sourceMappingURL inside a string literal), got '%s'", url)
		}
	})

	t.Run("Finds sourceMappingURL at end after code", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "after-code.js")
		content := "console.log('test');\n//# sourceMappingURL=valid.js.map"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "valid.js.map" {
			t.Errorf("expected 'valid.js.map', got '%s'", url)
		}
	})

	t.Run("Returns empty when comment contains */", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "with-close-comment.js")
		content := "console.log('test');\n//# sourceMappingURL=bad.js.map */\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "" {
			t.Errorf("expected empty string (comment contains */), got '%s'", url)
		}
	})

	t.Run("Handles whitespace-only lines", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "with-whitespace.js")
		content := "console.log('test');\n   \n//# sourceMappingURL=whitespace.js.map\n   "
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "whitespace.js.map" {
			t.Errorf("expected 'whitespace.js.map', got '%s'", url)
		}
	})

	t.Run("Processes lines in reverse order", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "reverse-order.js")
		// First comment should be ignored (in middle of file)
		// Second comment should be found (at end of file)
		content := "//# sourceMappingURL=wrong.js.map\nconsole.log('test');\n//# sourceMappingURL=correct.js.map\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "correct.js.map" {
			t.Errorf("expected 'correct.js.map' (from last comment), got '%s'", url)
		}
	})

	t.Run("Stops at non-comment content", func(t *testing.T) {
		bundlePath := filepath.Join(tempDir, "non-comment.js")
		// Should find nothing because there's code after the sourceMappingURL
		content := "//# sourceMappingURL=before-code.js.map\nconsole.log('test');\n"
		if err := os.WriteFile(bundlePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		url, err := upload.ExtractSourceMappingURL(bundlePath, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Scanning from the end, we hit "console.log" before the comment
		if url != "" {
			t.Errorf("expected empty string (code appears after sourceMappingURL when reading from end), got '%s'", url)
		}
	})
}

func TestResolveBundlePaths(t *testing.T) {
	logger := log.NewLoggerWrapper("debug")
	tempDir := t.TempDir()

	// Create various bundle files
	jsBundle := filepath.Join(tempDir, "app.js")
	if err := os.WriteFile(jsBundle, []byte("console.log('js');"), 0644); err != nil {
		t.Fatalf("failed to write js bundle: %v", err)
	}

	jsxBundle := filepath.Join(tempDir, "component.jsx")
	if err := os.WriteFile(jsxBundle, []byte("export const Comp = () => {};"), 0644); err != nil {
		t.Fatalf("failed to write jsx bundle: %v", err)
	}

	tsBundle := filepath.Join(tempDir, "types.ts")
	if err := os.WriteFile(tsBundle, []byte("const x: number = 1;"), 0644); err != nil {
		t.Fatalf("failed to write ts bundle: %v", err)
	}

	tsxBundle := filepath.Join(tempDir, "tsx-comp.tsx")
	if err := os.WriteFile(tsxBundle, []byte("export const TSXComp = () => <div />;"), 0644); err != nil {
		t.Fatalf("failed to write tsx bundle: %v", err)
	}

	// Should be ignored
	cssFile := filepath.Join(tempDir, "styles.css")
	if err := os.WriteFile(cssFile, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatalf("failed to write css file: %v", err)
	}

	// Should be ignored
	mjsBundle := filepath.Join(tempDir, "module.mjs")
	if err := os.WriteFile(mjsBundle, []byte("export const x = 1;"), 0644); err != nil {
		t.Fatalf("failed to write mjs bundle: %v", err)
	}

	t.Run("Finds all supported extensions", func(t *testing.T) {
		paths, err := upload.ResolveBundlePaths("", tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(paths) != 4 {
			t.Fatalf("expected 4 bundles, got %d", len(paths))
		}

		// Check that all expected files are found
		foundFiles := make(map[string]bool)
		for _, path := range paths {
			foundFiles[filepath.Base(path)] = true
		}

		expectedFiles := []string{"app.js", "component.jsx", "types.ts", "tsx-comp.tsx"}
		for _, expected := range expectedFiles {
			if !foundFiles[expected] {
				t.Errorf("expected to find %s, but it was not found", expected)
			}
		}

		// Check that unsupported files are not found
		if foundFiles["styles.css"] {
			t.Error("css file should not be included")
		}
		if foundFiles["module.mjs"] {
			t.Error("mjs file should not be included")
		}
	})

	t.Run("Returns explicit bundle path", func(t *testing.T) {
		paths, err := upload.ResolveBundlePaths(jsBundle, tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(paths) != 1 {
			t.Fatalf("expected 1 bundle, got %d", len(paths))
		}

		if paths[0] != jsBundle {
			t.Errorf("expected %s, got %s", jsBundle, paths[0])
		}
	})
}

// TestResolveSourceMapPaths_TIER3B_ViteHiddenSourcemaps tests the TIER 3B fallback mechanism
// for bundlers like Vite that use 'sourcemaps: hidden' configuration.
//
// Vite with hidden sourcemaps:
// - Creates .map files on disk
// - Does NOT add sourceMappingURL comments to bundle files
// - Standard ECMA-426 source map discovery fails
//
// This test verifies the fallback .map suffix matching works correctly.
func TestResolveSourceMapPaths_TIER3B_ViteHiddenSourcemaps(t *testing.T) {
	logger := log.NewLoggerWrapper("debug")

	t.Run("Finds source map by .map suffix when sourceMappingURL missing (Vite hidden mode)", func(t *testing.T) {
		tempDir := t.TempDir()
		
		// Simulate Vite hidden sourcemaps: bundle WITHOUT sourceMappingURL comment
		bundlePath := filepath.Join(tempDir, "app.js")
		bundleContent := "console.log('vite app');\n// No sourceMappingURL comment - hidden sourcemaps\n"
		if err := os.WriteFile(bundlePath, []byte(bundleContent), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		// Create corresponding .map file
		mapPath := filepath.Join(tempDir, "app.js.map")
		mapContent := `{"version":3,"sources":["src/main.ts"],"mappings":"test"}`
		if err := os.WriteFile(mapPath, []byte(mapContent), 0644); err != nil {
			t.Fatalf("failed to write source map: %v", err)
		}

		// Resolve source maps without explicit --source-map parameter (auto-discovery)
		bundles, err := upload.ResolveSourceMapPaths("", "", tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should find 1 bundle+map pair via TIER 3B fallback
		if len(bundles) != 1 {
			t.Fatalf("expected 1 source map bundle (TIER 3B fallback), got %d", len(bundles))
		}

		if bundles[0].BundlePath != bundlePath {
			t.Errorf("expected bundle path %s, got %s", bundlePath, bundles[0].BundlePath)
		}

		if bundles[0].SourceMapPath != mapPath {
			t.Errorf("expected map path %s, got %s", mapPath, bundles[0].SourceMapPath)
		}
	})

	t.Run("Prefers sourceMappingURL comment over .map suffix fallback", func(t *testing.T) {
		tempDir := t.TempDir()
		
		// Bundle WITH sourceMappingURL comment (priority over suffix matching)
		bundlePath := filepath.Join(tempDir, "priority.js")
		bundleContent := "console.log('test');\n//# sourceMappingURL=custom.js.map\n"
		if err := os.WriteFile(bundlePath, []byte(bundleContent), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		// Create the referenced source map
		customMapPath := filepath.Join(tempDir, "custom.js.map")
		if err := os.WriteFile(customMapPath, []byte("{}"), 0644); err != nil {
			t.Fatalf("failed to write custom map: %v", err)
		}

		// Also create a .map suffix file (should NOT be used)
		suffixMapPath := filepath.Join(tempDir, "priority.js.map")
		if err := os.WriteFile(suffixMapPath, []byte("{}"), 0644); err != nil {
			t.Fatalf("failed to write suffix map: %v", err)
		}

		bundles, err := upload.ResolveSourceMapPaths("", "", tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should use the sourceMappingURL reference, NOT the .map suffix file
		if len(bundles) != 1 {
			t.Fatalf("expected 1 source map bundle, got %d", len(bundles))
		}

		if bundles[0].SourceMapPath != customMapPath {
			t.Errorf("expected to use sourceMappingURL reference %s, but got %s", customMapPath, bundles[0].SourceMapPath)
		}
	})

	t.Run("Skips bundle when no sourceMappingURL and no .map suffix file exists", func(t *testing.T) {
		tempDir := t.TempDir()
		
		// Bundle WITHOUT sourceMappingURL comment
		orphanBundle := filepath.Join(tempDir, "orphan.js")
		if err := os.WriteFile(orphanBundle, []byte("console.log('orphan');"), 0644); err != nil {
			t.Fatalf("failed to write bundle: %v", err)
		}

		// NO .map file - should be skipped

		bundles, err := upload.ResolveSourceMapPaths("", "", tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should find nothing for this orphan bundle (neither comment nor suffix file)
		if len(bundles) != 0 {
			t.Fatalf("expected 0 source map bundles (orphan bundle without map file), got %d", len(bundles))
		}
	})

	t.Run("Works with React/TypeScript bundle names", func(t *testing.T) {
		tempDir := t.TempDir()
		
		// Realistic Vite React+TS bundle name: bundle-abc123.min.js
		complexBundlePath := filepath.Join(tempDir, "index-abc123.min.js")
		if err := os.WriteFile(complexBundlePath, []byte("var app={};"), 0644); err != nil {
			t.Fatalf("failed to write complex bundle: %v", err)
		}

		// Corresponding .map file with same naming
		complexMapPath := filepath.Join(tempDir, "index-abc123.min.js.map")
		if err := os.WriteFile(complexMapPath, []byte("{}"), 0644); err != nil {
			t.Fatalf("failed to write complex map: %v", err)
		}

		bundles, err := upload.ResolveSourceMapPaths("", "", tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(bundles) != 1 {
			t.Fatalf("expected 1 source map bundle, got %d", len(bundles))
		}

		if bundles[0].SourceMapPath != complexMapPath {
			t.Errorf("expected map path %s, got %s", complexMapPath, bundles[0].SourceMapPath)
		}
	})

	t.Run("TIER 3B fallback only used in auto-discovery mode (no explicit params)", func(t *testing.T) {
		tempDir := t.TempDir()
		
		// When explicit --source-map and --bundle are provided, use those (TIER 1)
		// When explicit --source-map only, find bundle by suffix (TIER 2)
		// When neither, auto-discover via sourceMappingURL (TIER 3A)
		// When neither + no sourceMappingURL, try .map suffix (TIER 3B)

		// For TIER 3B test: no explicit params = auto-discovery
		testBundle := filepath.Join(tempDir, "tier3b.js")
		if err := os.WriteFile(testBundle, []byte("var x=1;"), 0644); err != nil {
			t.Fatalf("failed to write test bundle: %v", err)
		}

		testMap := filepath.Join(tempDir, "tier3b.js.map")
		if err := os.WriteFile(testMap, []byte("{}"), 0644); err != nil {
			t.Fatalf("failed to write test map: %v", err)
		}

		// Call with no explicit sourceMapPath or bundlePath (auto-discovery mode)
		bundles, err := upload.ResolveSourceMapPaths("", "", tempDir, logger)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(bundles) != 1 {
			t.Fatalf("expected 1 bundle in TIER 3B fallback, got %d", len(bundles))
		}
	})
}
