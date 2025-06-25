package utils_testing

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"testing"
)

func TestValidateEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		apiKey   string
		expected string
	}{
		{
			name:     "Redirect upload.bugsnag.com when apiKey starts with 00000",
			endpoint: "https://upload.bugsnag.com",
			apiKey:   "0000012345",
			expected: "https://upload.insighthub.smartbear.com",
		},
		{
			name:     "Redirect upload.bugsnag.com/ when apiKey starts with 00000",
			endpoint: "https://upload.bugsnag.com/",
			apiKey:   "0000099999",
			expected: "https://upload.insighthub.smartbear.com",
		},
		{
			name:     "Do not redirect upload.bugsnag.com when apiKey does not start with 00000",
			endpoint: "https://upload.bugsnag.com",
			apiKey:   "1234567890",
			expected: "https://upload.bugsnag.com",
		},

		{
			name:     "Redirect build.bugsnag.com when apiKey starts with 00000",
			endpoint: "https://build.bugsnag.com",
			apiKey:   "0000054321",
			expected: "https://build.insighthub.smartbear.com",
		},
		{
			name:     "Redirect build.bugsnag.com/ when apiKey starts with 00000",
			endpoint: "https://build.bugsnag.com/",
			apiKey:   "0000088888",
			expected: "https://build.insighthub.smartbear.com",
		},
		{
			name:     "Do not redirect build.bugsnag.com when apiKey does not start with 00000",
			endpoint: "https://build.bugsnag.com",
			apiKey:   "ABCDE12345",
			expected: "https://build.bugsnag.com",
		},

		{
			name:     "Custom endpoint with matching apiKey is unchanged",
			endpoint: "https://custom.example.com",
			apiKey:   "0000011111",
			expected: "https://custom.example.com",
		},
		{
			name:     "Custom endpoint with non-matching apiKey is unchanged",
			endpoint: "https://custom.example.com",
			apiKey:   "ZZZZZ",
			expected: "https://custom.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateEndpoint(tt.endpoint, tt.apiKey)
			if result != tt.expected {
				t.Errorf("ValidateEndpoint(%q, %q) = %q; want %q", tt.endpoint, tt.apiKey, result, tt.expected)
			}
		})
	}
}
