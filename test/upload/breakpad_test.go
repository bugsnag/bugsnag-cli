package upload_testing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUploadBreakpadSymbolsQueryParams(t *testing.T) {
	t.Log("Testing getting query params for uploading breakpad symbols")

	apiKey := "1234567890ABCDEF1234567890ABCDEF"
	projectRoot := "/features/breakpad/fixtures/breakpad-symbols.sym"
	overwrite := true

	queryParams := fmt.Sprintf("?api_key=%s&overwrite=%t&project_root=%s",
		strings.ReplaceAll(apiKey, " ", "%20"),
		overwrite,
		strings.ReplaceAll(projectRoot, " ", "%20"),
	)

	expectedQueryParams := "?api_key=1234567890ABCDEF1234567890ABCDEF&overwrite=true&project_root=/features/breakpad/fixtures/breakpad-symbols.sym"

	assert.Equal(t, expectedQueryParams, queryParams)
}
