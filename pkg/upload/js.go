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
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// resolveProjectRoot determines the project root directory.
// If a project root is provided, it returns that value.
// Otherwise, it defaults to searching up the directory tree.
//
// Parameters:
// - projectRoot: A user-specified project root directory (can be empty).
// - path: A fallback path in case retrieving the working directory fails.
//
// Returns:
// - The resolved project root directory.
func resolveProjectRoot(projectRoot, path string) string {
	if projectRoot != "" {
		return projectRoot
	}

	// Get absolute path
	checkPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	// Workout how many segments to go up
	segments := len(strings.Split(checkPath, string(filepath.Separator)))

	for i := 0; i < segments; i++ {
		if hasProjectRootFile(checkPath) {
			return checkPath
		}
		parent := filepath.Dir(checkPath)
		if parent == checkPath { // Stop if we reach the root directory
			break
		}
		checkPath = parent
	}

	return checkPath
}

// HasProjectRootFile determines whether the specified directory contains
// any of the common project root files, indicating it is likely the root of a project.
//
// The function checks for the existence of the following files:
// - "package.json"
// - "yarn.lock"
// - "lerna.json"
//
// Parameters:
// - dir: The directory path to check.
//
// Returns:
// - bool: Returns true if any of the project root files are found in the directory; otherwise, false.
func hasProjectRootFile(dir string) bool {
	projectFiles := []string{"package.json", "yarn.lock", "lerna.json"}
	for _, file := range projectFiles {
		if utils.FileExists(filepath.Join(dir, file)) {
			return true
		}
	}
	return false
}

// resolveVersion attempts to parse the version from the package.json file
// if a version is not provided on the command line.
//
// Parameters:
// - versionName: version string from CLI (can be empty).
// - path: directory path to start searching.
//
// Returns:
// - resolved version string or empty if not found.
func resolveVersion(versionName string, path string, logger log.Logger) string {
	if versionName != "" {
		return versionName
	}
	logger.Debug(fmt.Sprintf("Attempting to automatically resolve the version starting from: %s", path))
	checkPath, err := filepath.Abs(path)
	if err != nil {
		logger.Warn(fmt.Sprintf("when resolving the version, unable to make an absolute path %s: %s", path, err))
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
			logger.Warn(fmt.Sprintf("when resolving the version, unable to open %s: %s", packageJson, err))
			return ""
		}
		byteValue, err := io.ReadAll(file)
		if err != nil {
			logger.Warn(fmt.Sprintf("when resolving the version, unable to read %s: %s", packageJson, err))
			return ""
		}
		var parsedPackageJson map[string]interface{}
		err = json.Unmarshal(byteValue, &parsedPackageJson)
		if err != nil {
			logger.Warn(fmt.Sprintf("when resolving the version, unable to parse %s: %s", packageJson, err))
			return ""
		}
		if parsedPackageJson["version"] == nil {
			logger.Warn(fmt.Sprintf("when resolving the version, the required version field wasn't found in in %s", packageJson))
			return ""
		}
		appVersion := parsedPackageJson["version"].(string)
		logger.Info(fmt.Sprintf("Using app version from %s: %s", packageJson, appVersion))
		return appVersion
	}
	return ""
}

// resolveSourceMapPaths attempts to find the source map(s) by walking the build directory
// if a source map path is not specified.
//
// Parameters:
// - sourceMapPath: user-specified source map path (can be empty).
// - outputPath: directory to scan for source maps.
//
// Returns:
// - list of source map file paths, or an error.
func resolveSourceMapPaths(sourceMapPath string, outputPath string, logger log.Logger) ([]string, error) {
	if sourceMapPath != "" {
		if utils.FileExists(sourceMapPath) {
			logger.Debug(fmt.Sprintf("Using user specified source map file %s", sourceMapPath))
			return []string{sourceMapPath}, nil
		} else {
			return []string{}, fmt.Errorf("unable to find specified source map file: %s", sourceMapPath)
		}
	}

	var sourceMapPaths []string
	err := filepath.WalkDir(outputPath, func(fullPath string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".map") {
			if !strings.HasSuffix(dirEntry.Name(), ".css.map") {
				sourceMapPaths = append(sourceMapPaths, fullPath)
			} else {
				logger.Debug(fmt.Sprintf("Skipping .css.map file %s", fullPath))
			}
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	if len(sourceMapPaths) == 0 {
		logger.Warn(fmt.Sprintf("No source maps found in: %s", outputPath))
	} else {
		logger.Info(fmt.Sprintf("Found source map(s): %s", strings.Join(sourceMapPaths, ", ")))
	}
	return sourceMapPaths, nil
}

// ReadSourceMap reads a JSON source map into memory.
//
// Parameters:
// - path: path to the source map file.
//
// Returns:
// - parsed JSON map of the source map contents or error.
func ReadSourceMap(path string, logger log.Logger) (map[string]interface{}, error) {
	logger.Info(fmt.Sprintf("Reading sourcemap %s", path))
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("cannot open sourcemap at %s: %s", path, err)
	}
	var sourceMapContents map[string]interface{}
	err = json.Unmarshal(data, &sourceMapContents)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("cannot unmarshal sourcemap at %s: %s", path, err)
	}
	return sourceMapContents, nil
}

// addSourcesContent adds the sourcesContent to a source map section if missing.
//
// Returns:
// - true if the source map was modified.
func addSourcesContent(section map[string]interface{}, sourceMapPath string, logger log.Logger) bool {
	untypedSources, hasSources := section["sources"]
	if !hasSources {
		logger.Warn(fmt.Sprintf("Cannot find the required sources field in the source map at %s", sourceMapPath))
		return false
	}
	sources, sourcesValid := untypedSources.([]interface{})
	if !sourcesValid {
		logger.Warn(fmt.Sprintf("The sources field exists but is not a list when trying to parse the source map at %s", sourceMapPath))
		return false
	}

	// Skip if the sources content is already valid
	untypedSourcesContent, hasSourcesContent := section["sourcesContent"]
	if hasSourcesContent {
		sourcesContent, sourcesContentValid := untypedSourcesContent.([]interface{})
		if sourcesContentValid && len(sourcesContent) == len(sources) {
			logger.Debug(fmt.Sprintf("SourcesContent is already populated for the source map at %s", sourceMapPath))
			return false
		}
	}

	var sourcesContent []*string
	for _, sourcePath := range sources {
		sourcePath, isString := sourcePath.(string)
		// Skip null values
		if !isString {
			sourcesContent = append(sourcesContent, nil)
			continue
		}
		// Handle the default webpack prefix in accordance with https://webpack.js.org/configuration/output/#outputdevtoolmodulefilenametemplate
		sourcePath, isWebpack := strings.CutPrefix(sourcePath, "webpack://")
		if isWebpack {
			// Remove the namespace
			firstSlash := strings.Index(sourcePath, "/")
			if firstSlash != -1 && firstSlash+1 < len(sourcePath) {
				sourcePath = sourcePath[firstSlash+1:]
			}

			// Skip virtual webpack files
			if strings.Contains(sourcePath, "webpack/") {
				sourcesContent = append(sourcesContent, nil)
				continue
			}

			// Remove any loaders
			questionMark := strings.LastIndex(sourcePath, "?")
			if questionMark-1 > 0 {
				sourcePath = sourcePath[:questionMark-1]
			}
		}
		if !filepath.IsAbs(sourcePath) {
			// Resolve the path relative to the source map
			sourcePath, _ = filepath.Abs(filepath.Join(filepath.Dir(sourceMapPath), sourcePath))
		}
		logger.Debug(fmt.Sprintf("Attempting to read the source %s.", sourcePath))
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			logger.Warn(fmt.Sprintf("Cannot read referenced source file '%s': %s", sourcePath, err))
			sourcesContent = append(sourcesContent, nil)
		} else {
			contentString := string(content)
			sourcesContent = append(sourcesContent, &contentString)
		}
	}

	section["sourcesContent"] = sourcesContent
	return true
}

// AddSources adds sourcesContent to the source map if missing.
//
// Returns:
// - true if the source map was modified.
func AddSources(sourceMapContents map[string]interface{}, sourceMapPath string, logger log.Logger) bool {
	// Sources may be in several sections. See https://bit.ly/sourcemap.
	if sections, exists := sourceMapContents["sections"]; exists {
		if sources, ok := sections.([]map[string]interface{}); ok {
			anyModified := false
			for _, section := range sources {
				modified := addSourcesContent(section, sourceMapPath, logger)
				anyModified = modified || anyModified
			}
			return anyModified
		}
	} else {
		return addSourcesContent(sourceMapContents, sourceMapPath, logger)
	}
	return false
}

// resolveBundlePath attempts to find the bundle path by changing the extension
// of the source map if bundle path is not specified.
//
// Parameters:
// - bundlePath: user-specified bundle path (can be empty).
// - sourceMapPath: path to the source map.
//
// Returns:
// - resolved bundle path or error.
func resolveBundlePath(bundlePath string, sourceMapPath string, logger log.Logger) (string, error) {
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
		logger.Info(fmt.Sprintf("Automatically using the bundle at path %s based on stripping the .map suffix.", withoutSuffix))
		return withoutSuffix, nil
	}
	return "", nil
}

// uploadSingleSourceMap uploads a single source map.
//
// Parameters:
// - sourceMapPath: path to source map file.
// - bundlePath: path to bundle file.
// - bundleUrl: URL for the bundle.
// - versionName: version string.
// - codeBundleId: optional code bundle id.
// - projectRoot: project root directory.
// - options: CLI options.
// - logger: logger instance.
//
// Returns:
// - error if upload fails.
func uploadSingleSourceMap(sourceMapPath string, bundlePath string, bundleUrl string, versionName string, codeBundleId string, projectRoot string, options options.CLI, logger log.Logger) error {
	sourceMapContents, err := ReadSourceMap(sourceMapPath, logger)
	if err != nil {
		return err
	}

	var sourceMapFile server.FileField

	sourceMapModified := AddSources(sourceMapContents, sourceMapPath, logger)
	if sourceMapModified {
		logger.Info(fmt.Sprintf("Added sources content to source map from %s", sourceMapPath))
		encodedSourceMap, err := json.Marshal(sourceMapContents)
		if err != nil {
			return fmt.Errorf("failed generate valid source map JSON with original sources added: %s", err.Error())
		}
		sourceMapFile = server.InMemoryFile{Path: sourceMapPath, Data: encodedSourceMap}
	} else {
		// Directly use the local file if the source map wasn't modified.
		logger.Info(fmt.Sprintf("Uploading unmodified source map from %s", sourceMapPath))
		sourceMapFile = server.LocalFile(sourceMapPath)
	}

	uploadOptions, err := utils.BuildJsUploadOptions(options.ApiKey, versionName, codeBundleId, bundleUrl, projectRoot, options.Upload.Js.Overwrite)

	if err != nil {
		return fmt.Errorf("failed to build upload options: %s", err.Error())
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = sourceMapFile
	fileFieldData["minifiedFile"] = server.LocalFile(bundlePath)

	err = server.ProcessFileRequest(options.ApiKey, "/sourcemap", uploadOptions, fileFieldData, sourceMapPath, options, logger)

	if err != nil {
		return fmt.Errorf("encountered error when uploading js sourcemap: %s", err.Error())
	}

	return nil
}

// ProcessJs uploads JS source maps based on CLI options.
//
// Parameters:
// - options: CLI options.
// - logger: logger instance.
//
// Returns:
// - error if processing fails.
func ProcessJs(options options.CLI, logger log.Logger) error {
	jsOptions := options.Upload.Js
	for _, path := range jsOptions.Path {

		outputPath := path

		// Set a default value for projectRoot if it's not defined
		jsOptions.ProjectRoot = resolveProjectRoot(jsOptions.ProjectRoot, path)
		logger.Debug(fmt.Sprintf("Using project root %s", jsOptions.ProjectRoot))

		jsOptions.VersionName = resolveVersion(jsOptions.VersionName, path, logger)

		// Check that the source map(s) exists and error out if it doesn't
		sourceMapPaths, err := resolveSourceMapPaths(jsOptions.SourceMap, outputPath, logger)
		if err != nil {
			return err
		}

		// Check that we now have a source map path
		if len(sourceMapPaths) == 0 {
			return fmt.Errorf("could not find a source map, please specify the path by using --source-map")
		}

		// Ensure that the correct one of --bundle-url and --base-url is specified
		isFile := utils.FileExists(jsOptions.SourceMap) || !utils.IsDir(outputPath)

		if isFile && jsOptions.BundleUrl == "" {
			return fmt.Errorf("`--bundle-url` must be set when uploading a file")
		}
		if isFile && jsOptions.BaseUrl != "" {
			return fmt.Errorf("`--base-url` must not be set when uploading a file")
		}
		if !isFile && jsOptions.BaseUrl == "" {
			return fmt.Errorf("`--base-url` must be set when uploading from a directory")
		}
		if !isFile && jsOptions.BundleUrl != "" {
			return fmt.Errorf("`--bundle-url` must not be set when uploading from a directory")
		}

		// Add a slash if it is not already on the end of the base URL
		if len(jsOptions.BaseUrl) > 0 && jsOptions.BaseUrl[len(jsOptions.BaseUrl)-1] != '/' {
			jsOptions.BaseUrl += "/"
		}

		for _, sourceMapPath := range sourceMapPaths {

			bundlePath, err := resolveBundlePath(jsOptions.Bundle, sourceMapPath, logger)
			if err != nil {
				return err
			}

			var bundleUrl string
			if jsOptions.BundleUrl != "" {
				bundleUrl = jsOptions.BundleUrl
			} else {
				// For directory uploads, add the relative path of the bundle to the base URL
				bundleUrl = jsOptions.BaseUrl + strings.TrimPrefix(strings.TrimPrefix(bundlePath, path), "/")
				logger.Debug(fmt.Sprintf("Generated URL %s using the base URL %s", bundleUrl, jsOptions.BaseUrl))
			}

			err = uploadSingleSourceMap(sourceMapPath, bundlePath, bundleUrl, jsOptions.VersionName, jsOptions.CodeBundleId, jsOptions.ProjectRoot, options, logger)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
