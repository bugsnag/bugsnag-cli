package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type JsOptions struct {
	VersionName string      `help:"The version of the app that the source map applies to. Defaults to the version in the package.json file (if found)."`
	BundleUrl   string      `help:"For single file uploads, the URL of the minified JavaScript file that the source map relates to. Asterisks can be used as a wildcard."`
	BaseUrl     string      `help:"For directory-based uploads, the URL of the base directory for the minified JavaScript files that the source maps relate to. The relative path is appended onto this for each file. Asterisks can be used as a wildcard."`
	SourceMap   string      `help:"Path to the source map file. This usually has the .min.js extension." type:"path"`
	Bundle      string      `help:"Path to the minified JavaScript file that the source map relates to. If this is not provided then the file will be obtained when an error event is received." type:"path"`
	ProjectRoot string      `help:"Path of the root directory on the file system where the source map was built. This will be stripped from the file name in the displayed stack traces." type:"path"`
	Path        utils.Paths `arg:"" name:"path" help:"Path to a directory of source maps and bundles to upload" type:"path" default:"."`
}

// Resolve the project if it isn't specified using the current working directory
func resolveProjectRoot(projectRoot string, path string) string {
	if projectRoot != "" {
		return projectRoot
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		return path
	}
	return workingDirectory
}

// Attempt to parse information from the package.json file if values aren't provided on the command line
func resolveVersion(versionName string, path string, logger log.Logger) string {
	if versionName != "" {
		return versionName
	}
	checkPath, err := filepath.Abs(path)
	if err != nil {
		logger.Warn(fmt.Sprintf("unable to make an absolute path %s: %s", path, err))
		return ""
	}
	// Walk up the folder structure as far as possible
	for filepath.Dir(checkPath) != checkPath {
		packageJson := filepath.Join(checkPath, "package.json")
		if !utils.FileExists(packageJson) {
			checkPath = filepath.Dir(checkPath)
			continue
		}
		file, err := os.Open(packageJson)
		if err != nil {
			logger.Warn(fmt.Sprintf("unable to open %s: %s", packageJson, err))
			return ""
		}
		byteValue, err := io.ReadAll(file)
		if err != nil {
			logger.Warn(fmt.Sprintf("unable to read %s: %s", packageJson, err))
			return ""
		}
		var parsedPackageJson map[string]interface{}
		err = json.Unmarshal(byteValue, &parsedPackageJson)
		if err != nil {
			logger.Warn(fmt.Sprintf("unable to parse %s: %s", packageJson, err))
			return ""
		}
		if parsedPackageJson["version"] == nil {
			logger.Warn(fmt.Sprintf("no version found in %s", packageJson))
			return ""
		}
		appVersion := parsedPackageJson["version"].(string)
		logger.Info(fmt.Sprintf("Using app version from %s: %s", packageJson, appVersion))
		return appVersion
	}
	return ""
}

// Attempt to find the source maps by walking the build directory if it is not passed in to the command line
func resolveSourceMapPaths(sourceMapPath string, outputPath string, logger log.Logger) ([]string, error) {
	if sourceMapPath != "" {
		if utils.FileExists(sourceMapPath) {
			return []string{sourceMapPath}, nil
		} else {
			return []string{}, fmt.Errorf("unable to find specified source map file: %s", sourceMapPath)
		}
	}

	var sourceMapPaths []string
	err := filepath.WalkDir(outputPath, func(fullPath string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".map") && !strings.HasSuffix(dirEntry.Name(), ".css.map") {
			sourceMapPaths = append(sourceMapPaths, fullPath)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	if len(sourceMapPaths) == 0 {
		logger.Debug(fmt.Sprintf("No main source map found in: %s", outputPath))
	} else {
		logger.Debug(fmt.Sprintf("Found source maps: %s", strings.Join(sourceMapPaths, ", ")))
	}
	return sourceMapPaths, nil
}

// Attempt to find the bundle path by changing the extension of the source map if the bundle path is not passed in to the command line
func resolveBundlePath(bundlePath string, sourceMapPath string) (string, error) {
	if bundlePath != "" {
		if utils.FileExists(bundlePath) {
			return bundlePath, nil
		} else {
			return "", fmt.Errorf("unable to find specified bundle: %s", bundlePath)
		}
	}

	withoutSuffix, found := strings.CutSuffix(sourceMapPath, ".map")
	if !found {
		return "", nil
	}

	if utils.FileExists(withoutSuffix) {
		return withoutSuffix, nil
	}
	return "", nil
}

// Upload a single source map
func uploadSingleSourceMap(
	bundleUrl string,
	baseUrl string,
	bundlePath string,
	sourceMapPath string,
	apiKey string,
	appVersion string,
	projectRoot string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) {

	bundlePath, err := resolveBundlePath(bundlePath, sourceMapPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	url := bundleUrl
	if baseUrl != "" {
		_, fileName := filepath.Split(bundlePath)
		if err != nil {
			logger.Warn(fmt.Sprintf("cannot make path relative to create URL, skipping: %v", err))
			return
		}
		url = baseUrl + fileName
	}

	uploadOptions, err := utils.BuildJsUploadOptions(apiKey, appVersion, url, projectRoot, overwrite)

	if err != nil {
		logger.Fatal(err.Error())
	}

	fileFieldData := make(map[string]string)
	fileFieldData["sourceMap"] = sourceMapPath
	if bundlePath != "" {
		fileFieldData["minifiedFile"] = bundlePath
	}

	err = server.ProcessFileRequest(endpoint+"/sourcemap", uploadOptions, fileFieldData, timeout, retries, sourceMapPath, dryRun, logger)

	if err != nil {
		logger.Fatal(err.Error())
	}
}

func ProcessJs(
	apiKey string,
	versionName string,
	bundleUrl string,
	baseUrl string,
	sourceMapPath string,
	bundlePath string,
	projectRoot string,
	Path utils.Paths,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	for _, path := range Path {

		outputPath := path

		// Set a default value for projectRoot if it's not defined
		projectRoot := resolveProjectRoot(projectRoot, path)

		appVersion := resolveVersion(versionName, path, logger)

		// Check that the source map(s) exists and error out if it doesn't
		sourceMapPaths, err := resolveSourceMapPaths(sourceMapPath, outputPath, logger)
		if err != nil {
			return err
		}

		// Check that we now have a source map path
		if len(sourceMapPaths) == 0 {
			return fmt.Errorf("could not find a source map, please specify the path by using --source-map")
		}

		// Ensure that the correct one of --bundle-url and --base-url is specified
		isFile := utils.FileExists(sourceMapPath) || !utils.IsDir(outputPath)

		if isFile && bundleUrl == "" {
			return fmt.Errorf("`--bundle-url` must be set when uploading a file")
		}
		if isFile && baseUrl != "" {
			return fmt.Errorf("`--base-url` must not be set when uploading a file")
		}
		if !isFile && baseUrl == "" {
			return fmt.Errorf("`--base-url` must be set when uploading from a directory")
		}
		if !isFile && bundleUrl != "" {
			return fmt.Errorf("`--bundle-url` must not be set when uploading from a directory")
		}

		// Add a slash if it is not already on the end of the base URL
		if len(baseUrl) > 0 && baseUrl[len(baseUrl)-1] != '/' {
			baseUrl += "/"
		}

		for _, sourceMapPath := range sourceMapPaths {

			uploadSingleSourceMap(bundleUrl, baseUrl, bundlePath, sourceMapPath, apiKey, appVersion, projectRoot, endpoint, timeout, retries, overwrite, dryRun, logger)
		}

	}

	return nil
}
