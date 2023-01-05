package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"syscall"
)

type AndroidNdkMapping struct {
	AndroidNdkRoot  string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifestPath string            `help:"(required) Path to directory or file to upload" type:"path"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot     string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}

func ProcessAndroidNDK(paths []string, androidNdkRoot string, appManifestPath string, configuration string, projectRoot string, buildUuid string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {

	// Check if we have project root
	if projectRoot == "" {
		return fmt.Errorf("`--project-root` missing from options")
	}

	androidNdkRoot, err := GetAndroidNDKRoot(androidNdkRoot)

	if err != nil {
		return err
	}

	log.Info("Using Android NDK Root: " + androidNdkRoot)

	return nil
}

// GetAndroidNDKRoot - Returns a valid Android NDK root path
func GetAndroidNDKRoot(path string) (string, error) {
	if path == "" {
		envValue, envPresent := os.LookupEnv("ANDROID_NDK_ROOT")
		if envPresent {
			path = envValue
		} else {
			return "", fmt.Errorf("unable to find ANDROID_NDK_ROOT")
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, fmt.Errorf(path + " does not exist on the system")
	}

	return path, nil
}
