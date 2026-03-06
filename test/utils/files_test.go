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
}
