package utils_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGzipCompress(t *testing.T) {
	t.Log("Testing compressing a given file")
	results, err := utils.GzipCompress("../testdata/android/android-mapping.txt")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "../testdata/android/android-mapping.gz", results, "File should be compressed")
}
