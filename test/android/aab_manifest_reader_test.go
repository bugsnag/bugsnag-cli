package android_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/stretchr/testify/assert"
)

func TestAabManifestReader(t *testing.T) {
	t.Log("Testing reading data from an AAB manifest XML")
	results, err := android.ReadAabManifest("../testdata/android/aab/AndroidManifest.xml")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, map[string]string(map[string]string{"apiKey": "your-api-key", "buildUuid": "19ce65f2-3a0f-434d-bbea-142c3ff23c48", "package": "com.example.bugsnag.android", "versionCode": "1", "versionName": "1.0"}), results, "The map data should match")
}
