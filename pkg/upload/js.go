package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// SourceMapBundle represents a paired bundle file and its corresponding source map.
type SourceMapBundle struct {
	BundlePath    string
	SourceMapPath string
}

// Precompiled regexp for matching sourceMappingURL comments.
var sourceMapURLRe = regexp.MustCompile(`^[@#]\s*sourceMappingURL=(\S*?)\s*$`)

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

// ExtractSourceMappingURL extracts the sourceMappingURL from a JavaScript bundle file.
// Implements the JavaScriptExtractSourceMapURL algorithm from TC39 ECMA-426 spec
// (section 11.1.2.1), using the "without parsing" approach.
//
// Parameters:
// - filePath: path to the bundle file.
// - logger: logger instance.
//
// Returns:
// - The source map URL/path if found, empty string if not found.
// - Error only if file cannot be read.
func ExtractSourceMappingURL(filePath string, logger log.Logger) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot read bundle file %s: %s", filePath, err)
	}

	source := string(data)

	// Split by line terminators as defined in ECMA-426: CR+LF, LF, CR, U+2028, U+2029
	lines := splitLines(source)

	// Process lines in reverse order
	for i := len(lines) - 1; i >= 0; i-- {
		lineStr := lines[i]
		line := []rune(lineStr)
		position := 0
		lineLength := len(line)

		// Scan through the line character by character
		for position < lineLength {
			first := line[position]

			// Check for // (single-line comment start)
			if first == '/' && position+1 < lineLength {
				position++
				second := line[position]

				if second == '/' {
					// Found //, rest of line is a comment
					position++
					comment := string(line[position:])

					// Validate comment doesn't contain quotes or */
					if strings.ContainsAny(comment, "\"'`") {
						return "", nil
					}
					if strings.Contains(comment, "*/") {
						return "", nil
					}

					// Try to match sourceMappingURL pattern
					sourceMapURL := matchSourceMapURL(comment)
					if sourceMapURL != "" {
						return sourceMapURL, nil
					}

					// Comment processed and did not contain a sourceMappingURL.
					// Per the "scan from end" algorithm, the presence of any trailing
					// non-matching line comment should terminate the search.
					return "", nil
				} else {
					// Found / but not //, this is invalid
					return "", nil
				}
			} else if isWhitespace(first) {
				// Skip whitespace
				position++
			} else {
				// Found non-whitespace, non-comment token - stop searching
				return "", nil
			}
		}
	}

	return "", nil
}

// splitLines splits a string by JavaScript line terminators as defined in ECMA-426.
// Line terminators: CR+LF (\r\n), LF (\n), CR (\r), U+2028 (Line Separator), U+2029 (Paragraph Separator)
func splitLines(source string) []string {
	// Replace all line terminators with \n, then split
	source = strings.ReplaceAll(source, "\r\n", "\n")
	source = strings.ReplaceAll(source, "\r", "\n")
	source = strings.ReplaceAll(source, "\u2028", "\n")
	source = strings.ReplaceAll(source, "\u2029", "\n")
	return strings.Split(source, "\n")
}

// isWhitespace checks if a rune is ECMAScript whitespace.
// Per ECMA-262: Space, Tab, Vertical Tab, Form Feed, NBSP, ZWNBSP, and other Unicode "Space_Separator"
func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\v', '\f', '\u00A0', '\uFEFF':
		return true
	default:
		// Check for Unicode Space_Separator category (Zs)
		return r == '\u1680' || (r >= '\u2000' && r <= '\u200A') || r == '\u202F' || r == '\u205F' || r == '\u3000'
	}
}

// matchSourceMapURL implements the MatchSourceMapURL algorithm from ECMA-426 spec.
// Pattern: ^[@#]\s*sourceMappingURL=(\S*?)\s*$
func matchSourceMapURL(comment string) string {
	// Pattern matches: [@#] followed by optional whitespace, then "sourceMappingURL=", then non-whitespace URL, then optional trailing whitespace
	matches := sourceMapURLRe.FindStringSubmatch(comment)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// ResolveBundlePaths finds JavaScript bundle files in the specified directory.
// If bundlePath is explicitly provided, it returns that single path.
//
// Parameters:
// - bundlePath: user-specified bundle path (can be empty).
// - outputPath: directory to scan for bundle files.
// - logger: logger instance.
//
// Returns:
// - list of bundle file paths, or an error.
func ResolveBundlePaths(bundlePath string, outputPath string, logger log.Logger) ([]string, error) {
	if bundlePath != "" {
		if utils.FileExists(bundlePath) {
			logger.Debug(fmt.Sprintf("Using user specified bundle file %s", bundlePath))
			return []string{bundlePath}, nil
		} else {
			return []string{}, fmt.Errorf("unable to find specified bundle file: %s", bundlePath)
		}
	}

	var bundlePaths []string
	supportedExtensions := []string{".js", ".jsx", ".ts", ".tsx"}

	err := filepath.WalkDir(outputPath, func(fullPath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dirEntry.IsDir() {
			// Skip walking node_modules directories entirely to avoid slow traversals.
			if dirEntry.Name() == "node_modules" {
				logger.Debug(fmt.Sprintf("Skipping node_modules directory: %s", fullPath))
				return fs.SkipDir
			}
			return nil
		}

		// Skip files in node_modules
		if strings.Contains(fullPath, "node_modules") {
			logger.Debug(fmt.Sprintf("Skipping bundle in node_modules: %s", fullPath))
			return nil
		}

		// Check if file has a supported extension
		for _, ext := range supportedExtensions {
			if strings.HasSuffix(dirEntry.Name(), ext) {
				bundlePaths = append(bundlePaths, fullPath)
				break
			} else {
				logger.Debug(fmt.Sprintf("Skipping non supported filetype: %s", fullPath))
			}
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	if len(bundlePaths) > 0 {
		logger.Debug(fmt.Sprintf("Found %d potential bundle file(s)", len(bundlePaths)))
	}

	return bundlePaths, nil
}

// ResolveSourceMapPaths attempts to find source map(s) by scanning bundle files
// for sourceMappingURL comments.
//
// Parameters:
// - sourceMapPath: user-specified source map path (can be empty).
// - bundlePath: user-specified bundle path (can be empty).
// - outputPath: directory to scan for bundle files.
// - logger: logger instance.
//
// Returns:
// - list of SourceMapBundle pairs, or an error.
func ResolveSourceMapPaths(sourceMapPath string, bundlePath string, outputPath string, logger log.Logger) ([]SourceMapBundle, error) {
	// If both source map and bundle are explicitly specified, return them as a pair
	if sourceMapPath != "" && bundlePath != "" {
		if !utils.FileExists(sourceMapPath) {
			return []SourceMapBundle{}, fmt.Errorf("unable to find specified source map file: %s", sourceMapPath)
		}
		if !utils.FileExists(bundlePath) {
			return []SourceMapBundle{}, fmt.Errorf("unable to find specified bundle file: %s", bundlePath)
		}
		logger.Debug(fmt.Sprintf("Using user specified source map %s and bundle %s", sourceMapPath, bundlePath))
		return []SourceMapBundle{{BundlePath: bundlePath, SourceMapPath: sourceMapPath}}, nil
	}

	// If only source map is specified, try to find the bundle by stripping .map suffix (legacy behavior)
	if sourceMapPath != "" {
		if !utils.FileExists(sourceMapPath) {
			return []SourceMapBundle{}, fmt.Errorf("unable to find specified source map file: %s", sourceMapPath)
		}
		logger.Debug(fmt.Sprintf("Using user specified source map file %s", sourceMapPath))
		
		// Try to find bundle by stripping .map suffix
		withoutSuffix, found := strings.CutSuffix(sourceMapPath, ".map")
		if found && utils.FileExists(withoutSuffix) {
			logger.Debug(fmt.Sprintf("Automatically using the bundle at path %s based on stripping the .map suffix.", withoutSuffix))
			return []SourceMapBundle{{BundlePath: withoutSuffix, SourceMapPath: sourceMapPath}}, nil
		}
		// If no bundle found, return empty bundle path (bundle is optional)
		return []SourceMapBundle{{BundlePath: "", SourceMapPath: sourceMapPath}}, nil
	}

	// Discover bundles and extract source map URLs from sourceMappingURL comments
	bundlePaths, err := ResolveBundlePaths(bundlePath, outputPath, logger)
	if err != nil {
		return []SourceMapBundle{}, err
	}

	var results []SourceMapBundle
	for _, bundleFile := range bundlePaths {
		logger.Debug(fmt.Sprintf("Attempting to locate sourcemap in bundle %s", bundleFile))
		sourceMappingURL, err := ExtractSourceMappingURL(bundleFile, logger)
		if err != nil {
			logger.Warn(fmt.Sprintf("Error reading bundle file %s: %s", bundleFile, err))
			continue
		}

		if sourceMappingURL == "" {
			logger.Debug(fmt.Sprintf("No sourceMappingURL found in %s", bundleFile))
			continue
		}

		// Check for data URLs (inline source maps)
		if strings.HasPrefix(sourceMappingURL, "data:") {
			logger.Debug(fmt.Sprintf("Skipping inline source map (data URL) in %s", bundleFile))
			continue
		}

		// Resolve relative path
		var sourceMapFile string
		if filepath.IsAbs(sourceMappingURL) {
			sourceMapFile = sourceMappingURL
		} else {
			sourceMapFile = filepath.Join(filepath.Dir(bundleFile), sourceMappingURL)
		}

		// Check if source map file exists
		if !utils.FileExists(sourceMapFile) {
			logger.Warn(fmt.Sprintf("Source map file %s referenced in %s does not exist", sourceMapFile, bundleFile))
			continue
		}

		logger.Info(fmt.Sprintf("Found source map %s for bundle %s", sourceMapFile, bundleFile))
		results = append(results, SourceMapBundle{
			BundlePath:    bundleFile,
			SourceMapPath: sourceMapFile,
		})
	}

	if len(results) > 0 {
		logger.Debug(fmt.Sprintf("Found %d source map(s)", len(results)))
	}

	return results, nil
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

	uploadOptions, err := utils.BuildJsUploadOptions(versionName, codeBundleId, bundleUrl, projectRoot, options.Upload.Js.Overwrite)

	if err != nil {
		return fmt.Errorf("failed to build upload options: %s", err.Error())
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = sourceMapFile
	fileFieldData["minifiedFile"] = server.LocalFile(bundlePath)

	err = server.ProcessFileRequest(
		options.ApiKey,
		"/sourcemap",
		uploadOptions,
		fileFieldData,
		sourceMapPath,
		options,
		logger,
	)

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

		// If the path is a .map file and no explicit --source-map is set,
		// treat the path as the source map itself
		if jsOptions.SourceMap == "" && strings.HasSuffix(path, ".map") && utils.FileExists(path) && !utils.IsDir(path) {
			jsOptions.SourceMap = path
			// For a direct .map file, use its directory as the output path
			outputPath = filepath.Dir(path)
		}

		// Set a default value for projectRoot if it's not defined
		jsOptions.ProjectRoot = resolveProjectRoot(jsOptions.ProjectRoot, path)
		logger.Debug(fmt.Sprintf("Using project root %s", jsOptions.ProjectRoot))

		jsOptions.VersionName = resolveVersion(jsOptions.VersionName, path, logger)

		// Resolve source map and bundle pairs
		sourceMapBundles, err := ResolveSourceMapPaths(jsOptions.SourceMap, jsOptions.Bundle, outputPath, logger)
		if err != nil {
			return err
		}

		// Check that we found at least one source map
		if len(sourceMapBundles) == 0 {
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

		for _, bundle := range sourceMapBundles {
			var bundleUrl string
			if jsOptions.BundleUrl != "" {
				bundleUrl = jsOptions.BundleUrl
			} else {
				// For directory uploads, add the relative path of the bundle to the base URL
				bundleUrl = jsOptions.BaseUrl + strings.TrimPrefix(strings.TrimPrefix(bundle.BundlePath, path), "/")
				logger.Debug(fmt.Sprintf("Generated URL %s using the base URL %s", bundleUrl, jsOptions.BaseUrl))
			}

			err = uploadSingleSourceMap(bundle.SourceMapPath, bundle.BundlePath, bundleUrl, jsOptions.VersionName, jsOptions.CodeBundleId, jsOptions.ProjectRoot, options, logger)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
