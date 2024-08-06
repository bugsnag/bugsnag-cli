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

// Attempt to find the source maps by walking the build directory if it is not passed in to the command line
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

// Reads a JSON source map into memory. Required in order to check if the sourcesContent exists.
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

// Based on the source map specification https://bit.ly/sourcemap. Returns if the source map was modified.
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
			sourcePath = path.Join(projectRoot, sourcePath)
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

// If the source map doesn't have sourcesContent then attempt to add it. Returns true if a modifcation was performed.
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

// Attempt to find the bundle path by changing the extension of the source map if the bundle path is not passed in to the command line
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

// Upload a single source map
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
		_, fileName := filepath.Split(jsOptions.Bundle)
		url = jsOptions.BaseUrl + fileName
		logger.Debug(fmt.Sprintf("Generated URL %s using the base URL %s", url, jsOptions.BaseUrl))
	}

	uploadOptions, err := utils.BuildJsUploadOptions(options.ApiKey, jsOptions.VersionName, url, jsOptions.ProjectRoot, options.Upload.Overwrite)

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
