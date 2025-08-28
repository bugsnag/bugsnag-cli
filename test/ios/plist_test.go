package ios_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPlistData(t *testing.T) {
	t.Run("empty path returns error", func(t *testing.T) {
		data, err := ios.GetPlistData("")
		assert.Nil(t, data)
		assert.Error(t, err)
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		data, err := ios.GetPlistData("does_not_exist.plist")
		assert.Nil(t, data)
		assert.Error(t, err)
	})

	t.Run("valid plist file", func(t *testing.T) {
		// Create a temporary plist file
		tmpDir := t.TempDir()
		plistPath := filepath.Join(tmpDir, "Info.plist")

		plistContent := `
		<?xml version="1.0" encoding="UTF-8"?>
		<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
		<plist version="1.0">
		<dict>
			<key>CreationDate</key>
			<date>2025-08-27T18:30:26Z</date>
			<key>CFBundleIdentifier</key>
			<string>com.example.app</string>
			<key>CFBundleShortVersionString</key>
			<string>1.2.3</string>
			<key>CFBundleVersion</key>
			<string>456</string>
			<key>bugsnag</key>
			<dict>
				<key>apiKey</key>
				<string>1234567890abcdef</string>
			</dict>
		</dict>
		</plist>`

		require.NoError(t, os.WriteFile(plistPath, []byte(plistContent), 0644))

		data, err := ios.GetPlistData(plistPath)
		require.NoError(t, err)

		assert.Equal(t, "com.example.app", data.BundleIdentifier)
		assert.Equal(t, "1.2.3", data.VersionName)
		assert.Equal(t, "456", data.BundleVersion)
		assert.Equal(t, "1234567890abcdef", data.BugsnagProjectDetails.ApiKey)
	})

	t.Run("malformed plist returns error", func(t *testing.T) {
		tmpDir := t.TempDir()
		plistPath := filepath.Join(tmpDir, "Invalid.plist")

		require.NoError(t, os.WriteFile(plistPath, []byte("not a plist"), 0644))

		data, err := ios.GetPlistData(plistPath)
		assert.Nil(t, data)
		assert.Error(t, err)
	})
}
