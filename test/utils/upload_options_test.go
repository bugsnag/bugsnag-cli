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
			projectRoot:  "/Users/vagrant/git",
			expectedRoot: "/Users/vagrant/git",
			expectDebug:  false,
		},
		"TC-01/2 single trailing slash normalized, debug shown": {
			projectRoot:  "/Users/vagrant/git/",
			expectedRoot: "/Users/vagrant/git",
			expectDebug:  true,
		},
		"TC-01/3 double trailing slash normalized, debug shown": {
			projectRoot:  "/Users/vagrant/git//",
			expectedRoot: "/Users/vagrant/git",
			expectDebug:  true,
		},
		"TC-01 double internal slash normalized, debug shown": {
			projectRoot:  "/Users//vagrant/git",
			expectedRoot: "/Users/vagrant/git",
			expectDebug:  true,
		},

		// TC-02: Relative paths resolved to absolute, debug shown
		"TC-02/1 dot-relative resolved to CWD, debug shown": {
			projectRoot:  ".",
			expectedRoot: cwd,
			expectDebug:  true,
		},
		"TC-02/2 parent-relative resolved, debug shown": {
			projectRoot:  "../",
			expectedRoot: strings.TrimSuffix(cwd, "/"+lastSegment(cwd)),
			expectDebug:  true,
		},

		// TC-03: Empty — debug shown (no projectRoot set)
		"TC-03/1 empty project root — no entry, debug shown": {
			projectRoot:  "",
			expectedRoot: "",
			expectDebug:  true,
		},

		// TC-04: Boundary values
		"TC-04/1 filesystem root / — debug shown": {
			projectRoot:  "/",
			expectedRoot: "/",
			expectDebug:  true,
		},
		"TC-04/2 path longer than 1024 characters — hard error": {
			projectRoot:   "/" + strings.Repeat("a", 1025),
			expectError:   true,
			errorContains: "exceeds the maximum allowed length",
		},

		// TC-05: Special characters
		"TC-05/1 path with spaces — no debug message": {
			projectRoot:  "/Users/dimple agarwal/Documents/GitHub/bugsnag-cocoa",
			expectedRoot: "/Users/dimple agarwal/Documents/GitHub/bugsnag-cocoa",
			expectDebug:  false,
		},
		"TC-05/2 path with unicode — no debug message": {
			projectRoot:  "/Users/dimple_agarwal/Documents/GitHub/bugsnag-cocoa\u00e9",
			expectedRoot: "/Users/dimple_agarwal/Documents/GitHub/bugsnag-cocoa\u00e9",
			expectDebug:  false,
		},

		// TC-07: Empty string
		"TC-07/1 empty string — debug shown": {
			projectRoot:  "",
			expectedRoot: "",
			expectDebug:  true,
		},

		// No leading slash: prepended and cleaned, debug shown
		"no leading slash — prepended to absolute, debug shown": {
			projectRoot:  "Users/rohan_dhiman",
			expectedRoot: "/Users/rohan_dhiman",
			expectDebug:  true,
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

