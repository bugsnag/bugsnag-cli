package utils_testing

import (
	"log"
	"os"
	"regexp"
	//"strconv"
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
	results := utils.IsDir(GetBasePath())

	assert.Equal(t, results, true, "This should be true")

	results = utils.IsDir(GetBasePath() + "/README.md")

	assert.Equal(t, results, false, "This should be false")
}

// TestBuildFileList - Tests the BuildFileList function
func TestBuildFileList(t *testing.T) {
	paths := []string{GetBasePath() + "/test/testdata"}

	results, err := utils.BuildFileList(paths)

	if err !=nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, results, []string{GetBasePath() + "/test/testdata/android-mapping.txt"}, "The files should be the same")
}

// TestFilePathWalkDir - Tests the FilePathWalkDir function
func TestFilePathWalkDir(t *testing.T) {
	results, err := utils.FilePathWalkDir(GetBasePath() + "/test/testdata")

	if err !=nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, results, []string{GetBasePath() + "/test/testdata/android-mapping.txt"}, "This should return a file")
}
