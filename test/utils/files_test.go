package utils_testing

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

//Pwd - Gets the current working directory
func GetTestDataDir() string {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	return path + "/../testdata"
}

// TestIsDir - Tests the IsDir function
func TestIsDir(t *testing.T) {
	got := utils.IsDir(GetTestDataDir())
	want := true

	if got != want {
		t.Errorf("got %q, wanted %q - %q", strconv.FormatBool(got), strconv.FormatBool(want), GetTestDataDir())
	}

	got = utils.IsDir(GetTestDataDir() + "/android-mapping.txt")
	want = false

	if got != want {
		t.Errorf("got %q, wanted %q", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

// TestBuildFileList - Tests the BuildFileList function
func TestBuildFileList(t *testing.T) {
	paths := []string{GetTestDataDir()}
	got, err := utils.BuildFileList(paths)

	if err !=nil {
		t.Errorf(err.Error())
	}

	want := []string{"/Users/josh.edney/repos/bugsnag-cli/test/testdata/android-mapping.txt"}

	if got[0] != want[0] {
		t.Errorf("got %q, want %q", got[0], want[0])
	}
}

// TestFilePathWalkDir - Tests the FilePathWalkDir function
func TestFilePathWalkDir(t *testing.T) {}
