package utils_testing

import (
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

//GetBasePath - Gets the current working directory
func GetBasePath() string {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	sampleRegexp := regexp.MustCompile(`/[^/]*/[^/]*$`)
	basePath := sampleRegexp.ReplaceAllString(path, "")

	return basePath
}

// TestIsDir - Tests the IsDir function
func TestIsDir(t *testing.T) {
	t.Log("Testing given path is a directory")
	results := utils.IsDir(GetBasePath())
	assert.Equal(t, results, true, "This should be true")

	t.Log("Testing given path is not a directory")
	results = utils.IsDir(GetBasePath() + "/README.md")
	assert.Equal(t, results, false, "This should be false")
}

// TestBuildFileList - Tests the BuildFileList function
func TestBuildFileList(t *testing.T) {
	t.Log("Testing building a file list from a given directory and file")
	paths := []string{GetBasePath() + "/test/testdata/android", GetBasePath() + "/README.md"}
	results, err := utils.BuildFileList(paths)

	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, results, []string{GetBasePath() + "/test/testdata/android/android-mapping.txt", GetBasePath() + "/README.md"}, "The files should be the same")

	t.Log("Testing building a file list from a single given file")
	paths = []string{GetBasePath() + "/test/testdata/android/android-mapping.txt"}
	results, err = utils.BuildFileList(paths)

	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, results, []string{GetBasePath() + "/test/testdata/android/android-mapping.txt"}, "The files should be the same")
}

// TestFilePathWalkDir - Tests the FilePathWalkDir function
func TestFilePathWalkDir(t *testing.T) {
	t.Log("Testing finding files within a given directory")
	results, err := utils.FilePathWalkDir(GetBasePath() + "/test/testdata/android")
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, results, []string{GetBasePath() + "/test/testdata/android/android-mapping.txt"}, "This should return a file")
}
