package utils_testing

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"testing"
)

func TestBuildEndpointUrl(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		port     int
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid URI without port, adds port",
			uri:      "http://localhost",
			port:     8080,
			expected: "http://localhost:8080",
			wantErr:  false,
		},
		{
			name:     "Valid URI with port, does not change",
			uri:      "http://localhost:3000",
			port:     8080,
			expected: "http://localhost:3000",
			wantErr:  false,
		},
		{
			name:     "Valid URI with no port and port=0",
			uri:      "http://127.0.0.1",
			port:     0,
			expected: "http://127.0.0.1",
			wantErr:  false,
		},
		{
			name:     "Invalid URI returns error",
			uri:      "http//invalid-url",
			port:     8080,
			expected: "http//invalid-url", // baseUrl.String() still returns what it could parse
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.BuildEndpointURL(tt.uri, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildEndpointUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
			if result != tt.expected {
				t.Errorf("BuildEndpointUrl() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDefaultUploadEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		apiKey   string
		expected string
	}{
		{
			name:     "Bugsnag upload endpoint with hub API key",
			endpoint: "https://upload.bugsnag.com",
			apiKey:   "00000abcde",
			expected: "https://upload.insighthub.smartbear.com",
		},
		{
			name:     "Bugsnag upload endpoint with regular API key",
			endpoint: "https://upload.bugsnag.com",
			apiKey:   "12345abcde",
			expected: "https://upload.bugsnag.com",
		},
		{
			name:     "Custom upload endpoint is unchanged",
			endpoint: "https://custom-upload.com",
			apiKey:   "00000abcde",
			expected: "https://custom-upload.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetDefaultUploadEndpoint(tt.endpoint, tt.apiKey)
			if result != tt.expected {
				t.Errorf("GetDefaultUploadEndpoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDefaultBuildEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		apiKey   string
		expected string
	}{
		{
			name:     "Bugsnag build endpoint with hub API key",
			endpoint: "https://build.bugsnag.com",
			apiKey:   "00000xyz",
			expected: "https://build.insighthub.smartbear.com",
		},
		{
			name:     "Bugsnag build endpoint with regular API key",
			endpoint: "https://build.bugsnag.com",
			apiKey:   "nonhubapikey",
			expected: "https://build.bugsnag.com",
		},
		{
			name:     "Custom build endpoint is unchanged",
			endpoint: "https://custom-build.com",
			apiKey:   "00000xyz",
			expected: "https://custom-build.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetDefaultBuildEndpoint(tt.endpoint, tt.apiKey)
			if result != tt.expected {
				t.Errorf("GetDefaultBuildEndpoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}
