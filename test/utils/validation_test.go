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
			name:     "Redirects when endpoint is upload.bugsnag.com and apiKey starts with 00000",
			endpoint: "https://upload.bugsnag.com",
			apiKey:   "0000012345",
			expected: "https://upload.insighthub.smartbear.com",
		},
		{
			name:     "Redirects when endpoint ends with slash and apiKey starts with 00000",
			endpoint: "https://upload.bugsnag.com/",
			apiKey:   "0000099999",
			expected: "https://upload.insighthub.smartbear.com",
		},
		{
			name:     "Does not redirect when endpoint is different",
			endpoint: "https://custom.endpoint.com",
			apiKey:   "0000012345",
			expected: "https://custom.endpoint.com",
		},
		{
			name:     "Does not redirect when apiKey does not start with 00000",
			endpoint: "https://upload.bugsnag.com",
			apiKey:   "1234500000",
			expected: "https://upload.bugsnag.com",
		},
		{
			name:     "Does not redirect when both endpoint and apiKey do not match criteria",
			endpoint: "https://example.com",
			apiKey:   "1234567890",
			expected: "https://example.com",
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
