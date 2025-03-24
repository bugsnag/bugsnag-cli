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
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// CleanPath securely normalizes a file path by removing any instances of "..",
// preventing directory traversal attacks.
//
// Parameters:
// - path: The input file path to be sanitized.
//
// Returns:
// - A cleaned and safe file path that does not contain directory traversal sequences.
func CleanPath(path string) string {
	// Normalize and clean the provided path
	cleaned := filepath.Clean(path)

	// Split the path into individual components
	parts := strings.Split(cleaned, string(filepath.Separator))
	var safeParts []string

	// Filter out any ".." components to prevent directory traversal
	for _, part := range parts {
		if part == ".." {
			continue
		}
		safeParts = append(safeParts, part)
	}

	// Reconstruct the sanitized path
	return filepath.Join(safeParts...)
}

// resolveProjectRoot determines the project root directory.
// If a project root is provided, it returns that value.
// Otherwise, it defaults to the current working directory.
//
// Parameters:
// - projectRoot: A user-specified project root directory (can be empty).
// - path: A fallback path in case retrieving the working directory fails.
//
// Returns:
// - The resolved project root directory.
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

// resolveVersion attempts to determine the application version by reading the package.json file
// if a version name is not provided via the command line.
//
// Parameters:
// - versionName: The explicitly provided version (if available).
// - path: The starting directory for searching package.json.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - The resolved application version as a string, or an empty string if resolution fails.
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

// resolveSourceMapPaths attempts to locate source map files by walking the build directory
// if a specific source map path is not provided via the command line.
//
// Parameters:
// - sourceMapPath: The explicitly provided source map file path (if available).
// - outputPath: The directory to search for source map files if none is provided.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - A slice of found source map file paths.
// - An error if file search encounters issues.
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

// ReadSourceMap reads a JSON source map file into memory.
// This is required to check if the 'sourcesContent' field exists.
//
// Parameters:
// - path: The file path to the source map.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - A map representing the parsed source map contents.
// - An error if the file cannot be opened or parsed.
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

// addSourcesContent modifies a source map by populating its 'sourcesContent' field.
// This is based on the source map specification: https://bit.ly/sourcemap.
//
// Parameters:
// - section: The source map section to modify.
// - sourceMapPath: Path to the source map file.
// - projectRoot: The root directory of the project.
// - logger: Logger instance for debugging and warnings.
//
// Returns:
// - A boolean indicating whether the source map was modified.
func addSourcesContent(section map[string]interface{}, sourceMapPath string, projectRoot string, logger log.Logger) bool {
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
		if !path.IsAbs(sourcePath) {
			sourcePath = CleanPath(sourcePath)
			sourcePath = filepath.Join(projectRoot, sourcePath)
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

// AddSources checks if a source map contains 'sourcesContent' and attempts to add it if missing.
// It processes both the top-level and sectioned source maps.
//
// Parameters:
// - sourceMapContents: Parsed JSON representation of the source map.
// - sourceMapPath: Path to the source map file.
// - projectRoot: Root directory of the project.
// - logger: Logger instance for logging messages.
//
// Returns:
// - true if any modifications were made, otherwise false.
func AddSources(sourceMapContents map[string]interface{}, sourceMapPath string, projectRoot string, logger log.Logger) bool {
	// Sources may be in several sections. See https://bit.ly/sourcemap.
	if sections, exists := sourceMapContents["sections"]; exists {
		if sources, ok := sections.([]map[string]interface{}); ok {
			anyModified := false
			for _, section := range sources {
				modified := addSourcesContent(section, sourceMapPath, projectRoot, logger)
				anyModified = modified || anyModified
			}
			return anyModified
		}
	} else {
		return addSourcesContent(sourceMapContents, sourceMapPath, projectRoot, logger)
	}
	return false
}

// resolveBundlePath attempts to determine the bundle file path by either using the provided path
// or by modifying the source map path if the bundle path is not explicitly set.
//
// Parameters:
// - bundlePath: The user-specified bundle path (if any).
// - sourceMapPath: The path to the source map file.
// - logger: Logger instance for logging messages.
//
// Returns:
// - The resolved bundle path if found, otherwise an error or empty string.
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

// uploadSingleSourceMap uploads a single JavaScript source map to the specified endpoint.
//
// Parameters:
// - options: CLI options provided by the user.
// - jsOptions: JavaScript-specific options, including source map details.
// - endpoint: The server endpoint for uploading source maps.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - An error if any part of the upload process fails, otherwise nil.
func uploadSingleSourceMap(options options.CLI, jsOptions options.Js, endpoint string, logger log.Logger) error {
	sourceMapContents, err := ReadSourceMap(jsOptions.SourceMap, logger)
	if err != nil {
		return err
	}

	var sourceMapFile server.FileField

	sourceMapModified := AddSources(sourceMapContents, jsOptions.SourceMap, jsOptions.ProjectRoot, logger)
	if sourceMapModified {
		logger.Info(fmt.Sprintf("Added sources content to source map from %s", jsOptions.SourceMap))
		encodedSourceMap, err := json.Marshal(sourceMapContents)
		if err != nil {
			return fmt.Errorf("failed generate valid source map JSON with original sources added: %s", err.Error())
		}
		sourceMapFile = server.InMemoryFile{Path: jsOptions.SourceMap, Data: encodedSourceMap}
	} else {
		// Directly use the local file if the source map wasn't modified.
		logger.Info(fmt.Sprintf("Uploading unmodified source map from %s", jsOptions.SourceMap))
		sourceMapFile = server.LocalFile(jsOptions.SourceMap)
	}

	jsOptions.Bundle, err = resolveBundlePath(jsOptions.Bundle, jsOptions.SourceMap, logger)
	if err != nil {
		return err
	}

	url := jsOptions.BundleUrl
	if jsOptions.BaseUrl != "" {
		fileName := strings.TrimPrefix(strings.TrimPrefix(jsOptions.Bundle, jsOptions.ProjectRoot), "/")

		newPath := strings.SplitN(fileName, "/", 2)
		if len(newPath) > 1 {
			fileName = newPath[1]
		}

		url = jsOptions.BaseUrl + fileName
		logger.Debug(fmt.Sprintf("Generated URL %s using the base URL %s", url, jsOptions.BaseUrl))
	}

	uploadOptions, err := utils.BuildJsUploadOptions(options.ApiKey, jsOptions.VersionName, jsOptions.CodeBundleId, url, jsOptions.ProjectRoot, options.Upload.Overwrite)

	if err != nil {
		return fmt.Errorf("failed to build upload options: %s", err.Error())
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = sourceMapFile
	if jsOptions.Bundle != "" {
		fileFieldData["minifiedFile"] = server.LocalFile(jsOptions.Bundle)
	}

	err = server.ProcessFileRequest(endpoint+"/sourcemap", uploadOptions, fileFieldData, jsOptions.SourceMap, options, logger)

	if err != nil {
		return fmt.Errorf("encountered error when uploading js sourcemap: %s", err.Error())
	}

	return nil
}

// ProcessJs handles the upload process for JavaScript source maps.
// It resolves paths, validates inputs, and ensures proper configuration
// before uploading each source map.
//
// Parameters:
// - options: CLI options containing configuration for the upload process.
// - endpoint: The API endpoint to which the source maps should be uploaded.
// - logger: Logger instance for logging messages.
//
// Returns:
// - error: An error if the process encounters an issue; otherwise, nil.
func ProcessJs(options options.CLI, endpoint string, logger log.Logger) error {
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
			jsOptions.SourceMap = sourceMapPath
			err := uploadSingleSourceMap(options, jsOptions, endpoint, logger)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
