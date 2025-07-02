package options_testing

import (
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"testing"
)

func TestMetadata_UnmarshalText(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		want      options.Metadata
		expectErr bool
	}{
		{
			name:  "Valid key-value pairs",
			input: []byte("key1=value1,key2=value2"),
			want:  options.Metadata{"key1": "value1", "key2": "value2"},
		},
		{
			name:      "Missing equal sign",
			input:     []byte("key1=value1,key2"),
			expectErr: true,
		},
		{
			name:  "Empty input",
			input: []byte(""),
			want:  options.Metadata{"": ""},
		},
		{
			name:  "Key with empty value",
			input: []byte("key1=,key2=value2"),
			want:  options.Metadata{"key1": "", "key2": "value2"},
		},
		{
			name:  "Value with equal sign",
			input: []byte("key1=val=ue1,key2=value2"),
			want:  options.Metadata{"key1": "val=ue1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m options.Metadata
			err := m.UnmarshalText(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(m) != len(tt.want) {
				t.Errorf("expected map length %d, got %d", len(tt.want), len(m))
			}
			for k, v := range tt.want {
				if m[k] != v {
					t.Errorf("for key %q, expected %q but got %q", k, v, m[k])
				}
			}
		})
	}
}
