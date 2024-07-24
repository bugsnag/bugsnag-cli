package upload_testing

import (
	"strings"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
)

func TestPopulateSourceMap(t *testing.T) {
	t.Log("Testing populating source map")
	logger := log.NewLoggerWrapper("debug")

	results, err := upload.ReadSourceMap("../testdata/js-nosources/dist/main.js.map", logger)
	if err != nil {
		t.Error(err)
	}

	modified := upload.AddSources(results, "../testdata/js-nosources", logger)
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
