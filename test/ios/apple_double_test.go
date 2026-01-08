package ios_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
)

func TestIsAppleDouble(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		content     []byte
		want        bool
		expectError bool
	}{
		{
			name:     "valid AppleDouble file",
			filename: "._file",
			content: []byte{
				0x00, 0x05, 0x16, 0x07,
				0x00, 0x02, 0x00, 0x00,
			},
			want:        true,
			expectError: false,
		},
		{
			name:        "regular file",
			filename:    "file.txt",
			content:     []byte("hello"),
			want:        false,
			expectError: false,
		},
		{
			name:        "empty file",
			filename:    "empty",
			content:     nil,
			want:        false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(path, tt.content, 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := ios.IsAppleDoubleMetaData(path)

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("IsAppleDouble(%s) = %v, want %v",
					tt.filename, got, tt.want)
			}
		})
	}
}
