package utils_testing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
)

// Tests expected common use case behaviour for processing <path> value
func TestProcessPathValue(t *testing.T) {

	// Get working dir
	currentDir, _ := os.Getwd()

	tt := map[string]struct {
		pathValue      string
		projectRoot    string
		expectedResult *ios.DsymUploadInfo
	}{
		"if <path> is set, is a normal directory and --project-root is not set, value of <path> is returned as-is": {
			pathValue:   "../testdata/ios/parent_root",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
			},
		},
		"if <path> is set, is a .xcodeproj directory and --project-root is not set, one directory up from <path> is returned": {
			pathValue:   "../testdata/ios/parent_root/MyTestApp.xcodeproj",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
			},
		},
		"if <path> is set, is a .xcworkspace directory and --project-root is not set, two directories up from <path> is returned": {
			pathValue:   "../testdata/ios/parent_root/MyTestApp.xcodeproj/project.xcworkspace",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
			},
		},
		"if <path> is set, is a normal directory and --project-root is set, --project-root takes precedence and it's value returned": {
			pathValue:   "../testdata/ios/parent_root",
			projectRoot: "../testdata/ios/alt_parent_root",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/alt_parent_root",
			},
		},
		"if <path> and --project-root are both unset, current working directory is returned": {
			pathValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: currentDir,
			},
		},
		"if <path> is a file then set the dsym path to it's value": {
			pathValue:   "../testdata/ios/MyTestApp",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				DsymPath: "../testdata/ios/MyTestApp",
			},
		},
		"if <path> is a .zip then set the dsym path to it's value": {
			pathValue:   "../testdata/ios/MyTestApp.zip",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				DsymPath: "../testdata/ios/MyTestApp.zip",
			},
		},
		"if <path> is a file then set the dsym path to it's value and if --project-root is set, use it's value for project root": {
			pathValue:   "../testdata/ios/MyTestApp",
			projectRoot: "../testdata/ios/parent_root",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
				DsymPath:    "../testdata/ios/MyTestApp",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actualResult, err := ios.ProcessPathValue(tc.pathValue, tc.projectRoot)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}
