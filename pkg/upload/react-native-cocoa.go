package upload

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeCocoa struct {
	AppVersion       string      `help:"The version of the application."`
	AppBundleVersion string      `help:"Bundle version for the application. (iOS only)"`
	Scheme           string      `help:"The name of the scheme to use when building the application."`
	SourceMap        string      `help:"Path to the source map file" type:"path"`
	Bundle           string      `help:"Path to the bundle file" type:"path"`
	Plist            string      `help:"Path to the Info.plist file" type:"path"`
	Xcworkspace      string      `help:"Path to the .xcworkspace file" type:"path"`
	CodeBundleID     string      `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev              bool        `help:"Indicates whether the application is a debug or release build"`
	ProjectRoot      string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Path             utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

type BuildSettings struct {
	ConfigurationBuildDir string `mapstructure:"CONFIGURATION_BUILD_DIR"`
	InfoPlistPath         string `mapstructure:"INFOPLIST_PATH"`
	BuiltProductsDir      string `mapstructure:"BUILT_PRODUCTS_DIR"`
}

type PlistData struct {
	AppVersion    string `mapstructure:"CFBundleShortVersionString"`
	BundleVersion string `mapstructure:"CFBundleVersion"`
	CodeBundleId  string `mapstructure:"CFBundleIdentifier"`
	ApiKey        string `mapstructure:"apiKey"`
}

func ProcessReactNativeCocoa(
	apiKey string,
	appVersion string,
	appBundleVersion string,
	scheme string,
	sourceMapPath string,
	bundlePath string,
	plistPath string,
	xcworkspacePath string,
	codeBundleId string,
	dev bool,
	projectRoot string,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {

	var err error
	var xcodeSchemes map[string]bool
	var buildSettings BuildSettings
	var plistData PlistData

	for _, path := range paths {

		// Set a default value for projectRoot if it's not defined
		if projectRoot == "" {
			projectRoot = path
		}

		// Set a default value for xcworkspacePath if it's not defined
		if xcworkspacePath == "" {
			iosPath := filepath.Join(path, "ios")
			xcworkspacePath, err = FindFolderWithSuffix(iosPath, ".xcworkspace")
			if err != nil {
				return err
			}
		}

		// Set a sourceMapPath if it's not defined and check that it exists before proceeding
		if sourceMapPath == "" {
			sourceMapPath = filepath.Join(projectRoot, "build", "sourcemaps", "main.jsbundle.map")
			if !utils.FileExists(sourceMapPath) {
				return errors.New("Could not find a suitable source map file, " +
					"please specify the path by using `--source-map`")
			}
		}

		// If the scheme is not defined, work out what the possible name is and retrieve all xcode schemes based on the xcworkspacePath
		if scheme == "" {
			possibleSchemeName := strings.TrimSuffix(filepath.Base(xcworkspacePath), ".xcworkspace")

			xcodeSchemes, err = GetXcodeSchemes(xcworkspacePath)
			if xcodeSchemes[possibleSchemeName] {
				// We can deduce that possibleSchemeName is the scheme name at this point, and can default to using it's value
				scheme = possibleSchemeName
				buildSettingsMap, err := GetXcodeBuildSettings(xcworkspacePath, possibleSchemeName)
				if err != nil {
					return err
				}
				err = mapstructure.Decode(buildSettingsMap, &buildSettings)

				// Set a default value for bundlePath if it's not defined and check that it exists before proceeding
				bundleFilePath := filepath.Join(buildSettings.ConfigurationBuildDir, "main.jsbundle")
				if !utils.FileExists(bundleFilePath) {
					return errors.New("Could not find a suitable bundle file, " +
						"please specify the path by using `--bundlePath`")
				}

				// Set a default value for plistPath if it's not defined
				if plistPath == "" {
					plistPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				}
				// TODO: Check Plistbuddy etc exists as well either within this func call or just before it
				plistDataMap, err := GetPlistData("/usr/libexec/Plistbuddy", plistPath)
				if err != nil {
					return err
				}
				err = mapstructure.Decode(plistDataMap, &plistData)

				// Set a default value for the relevant plist data needed for creating the upload options
				appBundleVersion = plistData.BundleVersion
				appVersion = plistData.AppVersion
				codeBundleId = plistData.CodeBundleId
				apiKey = plistData.ApiKey

			} else {
				return errors.New("Could not find a suitable scheme, " +
					"please specify the scheme by using `--scheme`")
			}

			if err != nil {
				return err
			}
		}

	}

	uploadOptions, err := utils.BuildReactNativeUploadOptions(apiKey, appVersion, appBundleVersion, codeBundleId, dev, projectRoot, overwrite, "cocoa")

	if err != nil {
		return err
	}

	fileFieldData := make(map[string]string)
	fileFieldData["sourceMap"] = sourceMapPath
	fileFieldData["bundle"] = bundlePath

	err = server.ProcessRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout, sourceMapPath, dryRun)

	if err != nil {
		return err
	} else {
		log.Success("Uploaded " + filepath.Base(sourceMapPath))
	}

	return nil
}

func GetPlistData(binary string, path string) (map[string]string, error) {
	cmd := exec.Command(binary, "-c", "Print :", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	variables := make(map[string]string)

	for _, line := range lines {
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			variables[key] = value
		}
	}

	return variables, nil
}

func GetXcodeSchemes(path string) (map[string]bool, error) {
	cmd := exec.Command("xcodebuild", "-workspace", path, "-list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	schemes := make(map[string]bool)
	lines := strings.Split(string(output), "\n")
	isParsingSchemes := false
	for _, line := range lines {
		if isParsingSchemes {
			if strings.TrimSpace(line) == "" {
				break
			}
			schemes[strings.TrimSpace(line)] = true
		} else if strings.Contains(line, "Schemes:") {
			isParsingSchemes = true
		}
	}
	return schemes, nil
}

func FindFolderWithSuffix(rootPath, targetSuffix string) (string, error) {
	var matchingFolder string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasSuffix(info.Name(), targetSuffix) {
			matchingFolder = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return matchingFolder, nil
}

func GetXcodeBuildSettings(path, schemeName string) (map[string]string, error) {
	cmd := exec.Command("xcodebuild", "-workspace", path, "-scheme", schemeName, "-showBuildSettings")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	buildSettings := make(map[string]string)
	var currentKey string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "    ") {
			buildSettings[currentKey] += line
		} else {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				buildSettings[key] = value
				currentKey = key
			}
		}
	}
	return buildSettings, nil
}
