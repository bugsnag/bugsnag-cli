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

// Resolve the project if it isn't specified using the currend working directory
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
func ResolveVersion(versionName string, projectRoot string, logger log.Logger) (string, error) {
	if versionName != "" {
		return versionName, nil
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
			logger.Warn(fmt.Sprintf("unable to open %s: %s", packageJson, err))
			return "", nil
		}
		byteValue, err := io.ReadAll(file)
		if err != nil {
			logger.Warn(fmt.Sprintf("unable to read %s: %s", packageJson, err))
			return "", nil
		}
		var parsedPackageJson map[string]interface{}
		err = json.Unmarshal(byteValue, &parsedPackageJson)
		if err != nil {
			logger.Warn(fmt.Sprintf("unable to parse %s: %s", packageJson, err))
			return "", nil
		}
		if parsedPackageJson["version"] == nil {
			logger.Warn(fmt.Sprintf("missing version in %s", packageJson))
			return "", nil
		}
		appVersion := parsedPackageJson["version"].(string)
		logger.Info(fmt.Sprintf("Using app version from %s: %s", packageJson, appVersion))
		return appVersion, nil
	}
	return "", nil
}

// Attempt to find the source maps by walking the build directory if it is not passed in to the command line
func ResolveSourceMapPaths(sourceMap string, outputPath string, projectRoot string, logger log.Logger) ([]string, error) {
	if sourceMap != "" {
		sourceMap := sourceMap
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
func ResolveBundlePath(
	bundle string, sourceMapPath string) (string, error) {
	if bundle != "" {
		if utils.FileExists(bundle) {
			return bundle, nil
		} else {
			return "", fmt.Errorf("unable to find specified bundle: %s", bundle)
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
	bundleUrl string,
	bundle string,
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

	bundlePath, err := ResolveBundlePath(bundle, sourceMapPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	uploadOptions, err := utils.BuildJsUploadOptions(apiKey, appVersion, bundleUrl, projectRoot, overwrite)

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
	sourceMap string,
	bundle string,
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

		// Dist is the default output path for webpack
		if utils.IsDir(filepath.Join(path, "dist")) {
			outputPath = filepath.Join(path, "dist")
		}

		// Set a default value for projectRoot if it's not defined
		projectRoot := resolveProjectRoot(projectRoot, path)

		appVersion, err := ResolveVersion(versionName, projectRoot, logger)
		if err != nil {
			return err
		}

		// Check that the source map(s) exists and error out if it doesn't
		sourceMapPaths, err := ResolveSourceMapPaths(sourceMap, outputPath, projectRoot, logger)
		if err != nil {
			return err
		}

		// Check that we now have a source map path
		if len(sourceMapPaths) == 0 {
			return fmt.Errorf("could not find a source map, please specify the path by using --source-map or SOURCEMAP_FILE environment variable")
		}

		for _, sourceMapPath := range sourceMapPaths {
			UploadSingleSourceMap(bundleUrl, bundle, sourceMapPath, apiKey, appVersion, projectRoot, endpoint, timeout, retries, overwrite, dryRun, logger)
		}

	}

	return nil
}
