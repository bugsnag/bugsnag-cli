package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

	log.Info("Locating ObjCopy within Android NDK Root")

	objCopyPath, err := BuildObjCopyPath(androidNdkRoot)

	if err != nil {
		return err
	}

	log.Info("Using ObjCopy located: " + objCopyPath)

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

// BuildObjCopyPath - Builds the path to the ObjCopy binary within the NDK root path
func BuildObjCopyPath(path string) (string, error) {
	ndkVersion, err := GetNdkVersion(path)
	if err != nil {
		return "", fmt.Errorf("unable to determine ndk version from path")
	}

	if ndkVersion < 24 {
		directoryPattern := filepath.Join(path, "/toolchains/x86_64-4.9/prebuilt/*/bin")
		directoryMatches, err := filepath.Glob(directoryPattern)
		if err != nil {
			return "", err
		}
		if directoryMatches == nil {
			return "", fmt.Errorf("Unable to find objcopy within ANDROID_NDK_ROOT: " + path)
		}

		if runtime.GOOS == "windows" {
			return filepath.Join(directoryMatches[0], "x86_64-linux-android-objcopy.exe"), nil
		}

		return filepath.Join(directoryMatches[0], "x86_64-linux-android-objcopy"), nil
	} else {
		directoryPattern := filepath.Join(path, "/toolchains/llvm/prebuilt/*/bin")
		directoryMatches, err := filepath.Glob(directoryPattern)
		if err != nil {
			return "", err
		}

		if directoryMatches == nil {
			return "", fmt.Errorf("Unable to find objcopy within ANDROID_NDK_ROOT: " + path)
		}

		if runtime.GOOS == "windows" {
			return filepath.Join(directoryMatches[0], "llvm-objcopy.exe"), nil
		}

		return filepath.Join(directoryMatches[0], "llvm-objcopy"), nil
	}
	return "", nil
}

// GetNdkVersion - Returns the major NDK version
func GetNdkVersion(path string) (int, error) {
	ndkVersion := strings.Split(filepath.Base(path), ".")
	ndkIntVersion, err := strconv.Atoi(ndkVersion[0])
	if err != nil {
		return 0, err
	}
	return ndkIntVersion, nil
}
