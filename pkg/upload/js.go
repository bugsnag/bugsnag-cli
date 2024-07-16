package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type JsOptions struct {
	VersionName string      `help:"The version of the app that the source map applies to. Defaults to the version in the package.json file (if found)."`
	BundleUrl   string      `help:"URL of the minified JavaScript file that the source map relates to. Asterisks can be used as a wildcard."`
	SourceMap   string      `help:"Path to the source map file. This usually has the .min.js extension." type:"path"`
	Bundle      string      `help:"Path to the minified JavaScript file that the source map relates to. If this is not provided then the file will be obtained when an error event is received." type:"path"`
	ProjectRoot string      `help:"Path of the root directory on the file system where the source map was built. This will be stripped from the file name in the displayed stack traces." type:"path"`
	Path        utils.Paths `arg:"" name:"path" help:"Path to a directory of source maps and bundles to upload" type:"path" default:"."`
}

// Resolve the project root by walking up from the path specified until a package.json is found
func resolveProjectRoot(jsOptions JsOptions, path string) (string, error) {
	if jsOptions.ProjectRoot != "" {
		return jsOptions.ProjectRoot, nil
	}
	checkPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("unable to make project root an absolute path %s", path)
	}
	// Walk up the folder structure as far as possible
	for filepath.Dir(checkPath) != checkPath {
		packageJson := filepath.Join(checkPath, "package.json")
		if !utils.FileExists(packageJson) {
			checkPath = filepath.Dir(checkPath)
			continue
		}
		return checkPath, nil
	}
	return path, nil
}

// Attempt to parse information from the package.json file if values aren't provided on the command line
func ResolveVersion(jsOptions JsOptions, projectRoot string, logger log.Logger) (string, error) {
	if jsOptions.VersionName != "" {
		return jsOptions.VersionName, nil
	}
	checkPath, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("unable to make project root an absolute path %s", projectRoot)
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
			return "", fmt.Errorf("unable to open package.json %s", packageJson)
		}
		byteValue, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("unable to read package.json %s", packageJson)
		}
		var parsedPackageJson map[string]interface{}
		err = json.Unmarshal(byteValue, &parsedPackageJson)
		if err != nil {
			return "", fmt.Errorf("unable to parse package.json %s %s", packageJson, err)
		}
		if parsedPackageJson["version"] == nil {
			return "", fmt.Errorf("package.json is missing the version %s", packageJson)
		}
		appVersion := parsedPackageJson["version"].(string)
		logger.Info(fmt.Sprintf("Using app version from package.json: %s", appVersion))
		return appVersion, nil
	}
	return "", fmt.Errorf("unable to locate package.json to resolve version in %s", projectRoot)
}

// Attempt to find the source maps by walking the build directory if it is not passed in to the command line
func ResolveSourceMapPaths(jsOptions JsOptions, outputPath string, projectRoot string, logger log.Logger) ([]string, error) {
	if jsOptions.SourceMap != "" {
		sourceMap := jsOptions.SourceMap
		if !path.IsAbs(sourceMap) {
			sourceMap = filepath.Join(projectRoot, sourceMap)
		}
		if utils.FileExists(sourceMap) {
			return []string{sourceMap}, nil
		} else {
			return []string{}, fmt.Errorf("unable to find specified source map file: %s", sourceMap)
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
func ResolveBundlePath(jsOptions JsOptions, sourceMapPath string) (string, error) {
	if jsOptions.Bundle != "" {
		if utils.FileExists(jsOptions.Bundle) {
			return jsOptions.Bundle, nil
		} else {
			return "", fmt.Errorf("unable to find specified bundle: %s", jsOptions.Bundle)
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
func UploadSingleSourceMap(
	jsOptions JsOptions,
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

	bundlePath, err := ResolveBundlePath(jsOptions, sourceMapPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	uploadOptions, err := utils.BuildJsUploadOptions(apiKey, appVersion, jsOptions.BundleUrl, projectRoot, overwrite)

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
	jsOptions JsOptions,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	for _, path := range jsOptions.Path {

		outputPath := path

		// Dist is the default output path for webpack
		if utils.IsDir(filepath.Join(path, "dist")) {
			outputPath = filepath.Join(path, "dist")
		}

		// Set a default value for projectRoot if it's not defined
		projectRoot, err := resolveProjectRoot(jsOptions, path)
		if err != nil {
			return err
		}

		appVersion, err := ResolveVersion(jsOptions, projectRoot, logger)
		if err != nil {
			return err
		}

		// Check that the source map(s) exists and error out if it doesn't
		sourceMapPaths, err := ResolveSourceMapPaths(jsOptions, outputPath, projectRoot, logger)
		if err != nil {
			return err
		}

		// Check that we now have a source map path
		if len(sourceMapPaths) == 0 {
			return fmt.Errorf("could not find a source map, please specify the path by using --source-map or SOURCEMAP_FILE environment variable")
		}

		for _, sourceMapPath := range sourceMapPaths {
			UploadSingleSourceMap(jsOptions, sourceMapPath, apiKey, appVersion, projectRoot, endpoint, timeout, retries, overwrite, dryRun, logger)
		}

	}

	return nil
}
