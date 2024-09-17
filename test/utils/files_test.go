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
