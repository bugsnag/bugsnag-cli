package utils_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// TestIsDir - Tests the IsDir function
func TestIsDir(t *testing.T) {
	t.Log("Testing given path is a directory")
	results := utils.IsDir("../../")
	assert.Equal(t, results, true, "Base path should be a directory")

	t.Log("Testing given path is not a directory")
	results = utils.IsDir("../../README.md")
	assert.Equal(t, results, false, "A regular file should not be a directory")
}

// TestBuildFileList - Tests the BuildFileList function
func TestBuildFileList(t *testing.T) {
	t.Log("Testing building a file list from a given directory and file")
	paths := []string{"../testdata/android/variants", "../../README.md"}
	results, err := utils.BuildFileList(paths)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	assert.Equal(t, results, []string{"../testdata/android/variants/debug/.gitkeep", "../testdata/android/variants/release/.gitkeep", "../../README.md"}, "The files should be the same")

	t.Log("Testing building a file list from a single given file")
	paths = []string{"../testdata/android/android-mapping.txt"}
	results, err = utils.BuildFileList(paths)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	assert.Equal(t, results, []string{"../testdata/android/android-mapping.txt"}, "The files should be the same")
}

// TestFilePathWalkDir - Tests the FilePathWalkDir function
func TestFilePathWalkDir(t *testing.T) {
	t.Log("Testing finding files within a given directory")
	results, err := utils.FilePathWalkDir("../testdata/android/variants")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	assert.Equal(t, results, []string{"../testdata/android/variants/debug/.gitkeep", "../testdata/android/variants/release/.gitkeep"}, "This should return a file")
}

// TestIsFileExcluded - Tests the IsFileExcluded function
func TestIsFileExcluded(t *testing.T) {
	t.Run("Matches wildcard extension pattern", func(t *testing.T) {
		excluded := utils.IsFileExcluded("path/to/file.map", []string{"*.map"})
		assert.True(t, excluded, "Should exclude files with .map extension")

		excluded = utils.IsFileExcluded("file.js.map", []string{"*.map"})
		assert.True(t, excluded, "Should exclude files ending with .map")
	})

	t.Run("Matches wildcard path pattern with basename", func(t *testing.T) {
		// filepath.Match works on basenames, so node_modules/* will match files in node_modules
		// but the actual match happens via substring matching in the implementation
		excluded := utils.IsFileExcluded("node_modules/package/file.js", []string{"node_modules"})
		assert.True(t, excluded, "Should exclude files with node_modules in path via substring")

		excluded = utils.IsFileExcluded("src/temp/file.js", []string{"temp"})
		assert.True(t, excluded, "Should exclude files with temp in path via substring")
	})

	t.Run("Matches exact filename", func(t *testing.T) {
		excluded := utils.IsFileExcluded("path/to/test.map", []string{"test.map"})
		assert.True(t, excluded, "Should exclude exact filename match")

		excluded = utils.IsFileExcluded("test.map", []string{"test.map"})
		assert.True(t, excluded, "Should exclude exact filename in current dir")
	})

	t.Run("Matches substring path", func(t *testing.T) {
		excluded := utils.IsFileExcluded("src/node_modules/lib/file.js", []string{"node_modules"})
		assert.True(t, excluded, "Should exclude files with path containing substring")

		excluded = utils.IsFileExcluded("dist/vendor/bundle.js", []string{"vendor"})
		assert.True(t, excluded, "Should exclude files with vendor in path")
	})

	t.Run("Does not match when pattern doesn't apply", func(t *testing.T) {
		excluded := utils.IsFileExcluded("src/main.js", []string{"*.map"})
		assert.False(t, excluded, "Should not exclude .js file with .map pattern")

		excluded = utils.IsFileExcluded("src/components/file.js", []string{"node_modules"})
		assert.False(t, excluded, "Should not exclude files without matching substring")
	})

	t.Run("Handles multiple patterns", func(t *testing.T) {
		patterns := []string{"*.map", "*.log", "node_modules"}

		excluded := utils.IsFileExcluded("file.map", patterns)
		assert.True(t, excluded, "Should match first pattern")

		excluded = utils.IsFileExcluded("debug.log", patterns)
		assert.True(t, excluded, "Should match second pattern")

		excluded = utils.IsFileExcluded("node_modules/lib/file.js", patterns)
		assert.True(t, excluded, "Should match third pattern")

		excluded = utils.IsFileExcluded("src/main.js", patterns)
		assert.False(t, excluded, "Should not match any pattern")
	})

	t.Run("Handles empty patterns", func(t *testing.T) {
		excluded := utils.IsFileExcluded("any/file.js", []string{})
		assert.False(t, excluded, "Should not exclude with no patterns")

		excluded = utils.IsFileExcluded("any/file.js", nil)
		assert.False(t, excluded, "Should not exclude with nil patterns")
	})

	t.Run("Handles complex wildcard patterns", func(t *testing.T) {
		excluded := utils.IsFileExcluded("test.js.map", []string{"*.js.map"})
		assert.True(t, excluded, "Should match .js.map extension")

		excluded = utils.IsFileExcluded("bundle-v1.2.3.js", []string{"bundle-*.js"})
		assert.True(t, excluded, "Should match bundle with version pattern")
	})

	t.Run("Handles directory path patterns", func(t *testing.T) {
		// Substring matching for directory paths
		excluded := utils.IsFileExcluded("build/dist/main.js", []string{"build"})
		assert.True(t, excluded, "Should match files with build in path")

		excluded = utils.IsFileExcluded("src/build/main.js", []string{"build"})
		assert.True(t, excluded, "Should match build as substring in path")

		excluded = utils.IsFileExcluded("src/main.js", []string{"build"})
		assert.False(t, excluded, "Should not match when build is not in path")
	})

	t.Run("Supports ** recursive globbing for directories", func(t *testing.T) {
		// node_modules/** should match all files under node_modules at the root level
		excluded := utils.IsFileExcluded("node_modules/package/file.js", []string{"node_modules/**"})
		assert.True(t, excluded, "Should exclude files in node_modules with ** pattern")

		excluded = utils.IsFileExcluded("node_modules/package/lib/deep/file.js", []string{"node_modules/**"})
		assert.True(t, excluded, "Should exclude deeply nested files in node_modules")

		// node_modules/** only matches if node_modules is at the start of the path
		excluded = utils.IsFileExcluded("src/node_modules/package/file.js", []string{"node_modules/**"})
		assert.False(t, excluded, "node_modules/** pattern only matches at path start")

		excluded = utils.IsFileExcluded("src/components/file.js", []string{"node_modules/**"})
		assert.False(t, excluded, "Should not exclude files outside node_modules")

		// To match node_modules at any level, use **/node_modules/**
		excluded = utils.IsFileExcluded("src/node_modules/package/file.js", []string{"**/node_modules/**"})
		assert.True(t, excluded, "Should exclude node_modules at any path level with **/node_modules/**")

		excluded = utils.IsFileExcluded("vendor/libs/node_modules/pkg/index.js", []string{"**/node_modules/**"})
		assert.True(t, excluded, "Should exclude node_modules deeply nested with ** pattern")
	})

	t.Run("Supports ** recursive globbing with wildcards", func(t *testing.T) {
		// **/*.map should match all .map files anywhere in the tree
		excluded := utils.IsFileExcluded("app.js.map", []string{"**/*.map"})
		assert.True(t, excluded, "Should match .map files in root")

		excluded = utils.IsFileExcluded("src/components/app.js.map", []string{"**/*.map"})
		assert.True(t, excluded, "Should match .map files in nested directories")

		excluded = utils.IsFileExcluded("build/dist/vendor/bundle.min.js.map", []string{"**/*.map"})
		assert.True(t, excluded, "Should match .map files deeply nested")

		excluded = utils.IsFileExcluded("src/app.js", []string{"**/*.map"})
		assert.False(t, excluded, "Should not match non-.map files")
	})

	t.Run("Supports ** in middle of path pattern", func(t *testing.T) {
		// **/temp/** should match any files in temp directories at any level
		excluded := utils.IsFileExcluded("temp/file.js", []string{"**/temp/**"})
		assert.True(t, excluded, "Should match files in root temp directory")

		excluded = utils.IsFileExcluded("src/temp/cache/file.js", []string{"**/temp/**"})
		assert.True(t, excluded, "Should match files in nested temp directory")

		excluded = utils.IsFileExcluded("build/output/temp/intermediate/file.js", []string{"**/temp/**"})
		assert.True(t, excluded, "Should match files in deeply nested temp directories")

		excluded = utils.IsFileExcluded("src/templates/file.js", []string{"**/temp/**"})
		assert.False(t, excluded, "Should not match files outside temp directories")
	})

	t.Run("Supports specific path with ** globbing", func(t *testing.T) {
		// src/**/*.test.js should match test files in src and subdirectories
		excluded := utils.IsFileExcluded("src/app.test.js", []string{"src/**/*.test.js"})
		assert.True(t, excluded, "Should match test files in src")

		excluded = utils.IsFileExcluded("src/components/button.test.js", []string{"src/**/*.test.js"})
		assert.True(t, excluded, "Should match test files in src subdirectories")

		excluded = utils.IsFileExcluded("src/utils/helpers/format.test.js", []string{"src/**/*.test.js"})
		assert.True(t, excluded, "Should match test files deeply nested in src")

		excluded = utils.IsFileExcluded("test/unit/app.test.js", []string{"src/**/*.test.js"})
		assert.False(t, excluded, "Should not match test files outside src")

		excluded = utils.IsFileExcluded("src/app.js", []string{"src/**/*.test.js"})
		assert.False(t, excluded, "Should not match non-test files in src")
	})

	t.Run("Combines ** patterns with other patterns", func(t *testing.T) {
		patterns := []string{"**/*.map", "node_modules/**", "**/dist/**", "*.log"}

		excluded := utils.IsFileExcluded("src/app.js.map", patterns)
		assert.True(t, excluded, "Should match .map pattern")

		excluded = utils.IsFileExcluded("node_modules/lib/index.js", patterns)
		assert.True(t, excluded, "Should match node_modules pattern")

		excluded = utils.IsFileExcluded("build/dist/bundle.js", patterns)
		assert.True(t, excluded, "Should match dist pattern")

		excluded = utils.IsFileExcluded("debug.log", patterns)
		assert.True(t, excluded, "Should match .log pattern")

		excluded = utils.IsFileExcluded("src/components/app.js", patterns)
		assert.False(t, excluded, "Should not match any pattern")
	})

	t.Run("Handles edge cases with ** patterns", func(t *testing.T) {
		// Test various edge cases
		excluded := utils.IsFileExcluded("file.js", []string{"**"})
		assert.True(t, excluded, "** should match everything")

		excluded = utils.IsFileExcluded("src/deep/path/file.js", []string{"**"})
		assert.True(t, excluded, "** should match any depth")

		excluded = utils.IsFileExcluded("a/b/c/d/file.js", []string{"**/b/**"})
		assert.True(t, excluded, "Should match with b directory in path")

		excluded = utils.IsFileExcluded("vendor/node_modules/pkg/index.js", []string{"**/node_modules/**"})
		assert.True(t, excluded, "Should match node_modules at any level")
	})
}
