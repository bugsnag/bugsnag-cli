package utils_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func TestEndpointBuilding(t *testing.T) {
	got := utils.BuildEndpointUrl("", 0)
	want := "https://upload.bugsnag.com"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = utils.BuildEndpointUrl("https://localhost", 0)
	want = "https://localhost"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = utils.BuildEndpointUrl("", 8443)
	want = "https://upload.bugsnag.com:8443"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = utils.BuildEndpointUrl("https://localhost", 8443)
	want = "https://localhost:8443"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
