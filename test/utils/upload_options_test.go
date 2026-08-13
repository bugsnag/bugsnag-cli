package utils_testing

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// TestBuildDsymUploadOptions_ProjectRoot covers the normalizeProjectRoot logic
// mapped to test case scenarios from the dSYM upload bug report.
// All messages are emitted at [DEBUG] level (visible only with --verbose).
func TestBuildDsymUploadOptions_ProjectRoot(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	tt := map[string]struct {
		projectRoot   string
		expectedRoot  string
		expectError   bool
		errorContains string
	}{
		// TC-01: Absolute path variants
		"TC-01/1 clean absolute path — no debug message": {
			projectRoot:  "/path/to/project",
			expectedRoot: "/path/to/project",
		},
		"TC-01/2 single trailing slash normalized": {
			projectRoot:  "/path/to/project/",
			expectedRoot: "/path/to/project",
		},
		"TC-01/3 double trailing slash normalized": {
			projectRoot:  "/path/to/project//",
			expectedRoot: "/path/to/project",
		},
		"TC-01 double internal slash normalized": {
			projectRoot:  "/path//to/project",
			expectedRoot: "/path/to/project",
		},

		// TC-02: Relative paths resolved to absolute
		"TC-02/1 dot-relative resolved to CWD": {
			projectRoot:  ".",
			expectedRoot: cwd,
		},
		"TC-02/2 parent-relative resolved": {
			projectRoot:  "../",
			expectedRoot: strings.TrimSuffix(cwd, "/"+lastSegment(cwd)),
		},

		// TC-03: Empty — no projectRoot set
		"TC-03/1 empty project root — no entry": {
			projectRoot:  "",
			expectedRoot: "",
		},

		// TC-04: Boundary values
		"TC-04/1 filesystem root /": {
			projectRoot:  "/",
			expectedRoot: "/",
		},
		"TC-04/2 path longer than 1024 characters — hard error": {
			projectRoot:   "/" + strings.Repeat("a", 1025),
			expectError:   true,
			errorContains: "exceeds the maximum allowed length",
		},

		// TC-05: Special characters
		"TC-05/1 path with spaces": {
			projectRoot:  "/path/to/my project/repo",
			expectedRoot: "/path/to/my project/repo",
		},
		"TC-05/2 path with unicode characters": {
			projectRoot:  "/path/to/project-\u00e9",
			expectedRoot: "/path/to/project-\u00e9",
		},

		// TC-07: Empty string
		"TC-07/1 empty string — no project root set": {
			projectRoot:  "",
			expectedRoot: "",
		},

		// No leading slash: prepended and cleaned
		"no leading slash — prepended to absolute": {
			projectRoot:  "path/to/project",
			expectedRoot: "/path/to/project",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			opts, err := utils.BuildDsymUploadOptions(tc.projectRoot)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				return
			}

			require.NoError(t, err)

			if tc.expectedRoot == "" {
				assert.NotContains(t, opts, "projectRoot", "projectRoot should not be set for empty input")
			} else {
				assert.Equal(t, tc.expectedRoot, opts["projectRoot"])
			}
		})
	}
}

// lastSegment returns the last path component of a slash-separated path.
func lastSegment(path string) string {
	path = strings.TrimRight(path, "/")
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return path
	}
	return path[idx+1:]
}
