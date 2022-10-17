package utils_testing

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
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
	got := utils.IsDir(GetBasePath())
	want := true

	if got != want {
		t.Errorf("got %q, wanted %q - %q", strconv.FormatBool(got), strconv.FormatBool(want), GetBasePath())
	}

	got = utils.IsDir(GetBasePath() + "/README.md")
	want = false

	if got != want {
		t.Errorf("got %q, wanted %q", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

// TestBuildFileList - Tests the BuildFileList function
func TestBuildFileList(t *testing.T) {
	paths := []string{GetBasePath() + "/test/testdata"}
	got, err := utils.BuildFileList(paths)

	if err !=nil {
		t.Errorf(err.Error())
	}

	want := []string{GetBasePath() + "/test/testdata/android-mapping.txt"}

	if got[0] != want[0] {
		t.Errorf("got %q, want %q", got[0], want[0])
	}
}

// TestFilePathWalkDir - Tests the FilePathWalkDir function
func TestFilePathWalkDir(t *testing.T) {}
