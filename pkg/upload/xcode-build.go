package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"os"
	"path/filepath"
)

func ProcessXcodeBuild(opts options.CLI, endpoint string, logger log.Logger) error {
	dsyms, plistPath, tempDirs, err := FindDsymsAndSettings(opts, logger)
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()
	if err != nil {
		return err
	}
	return UploadDsyms(dsyms, plistPath, endpoint, opts, logger)
}

func FindDsymsAndSettings(opts options.CLI, logger log.Logger) ([]*ios.DwarfInfo, string, []string, error) {
	xcodeBuildOptions := opts.Upload.XcodeBuild
	var (
		buildSettings *ios.XcodeBuildSettings
		dsymPath      string
		tempDirs      []string
		dwarfInfo     []*ios.DwarfInfo
		tempDir       string
		err           error
	)

	xcodeProjPath := string(xcodeBuildOptions.Shared.XcodeProject)
	plistPath := string(xcodeBuildOptions.Shared.Plist)

	for _, path := range xcodeBuildOptions.Path {
		if filepath.Ext(path) == ".xcarchive" {
			logger.Warn(fmt.Sprintf("The specified path %s is an Xcode archive. Please use the `xcode-archive` command instead as this functionality will be deprecated in future releases.", path))
		}

		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			if xcodeProjPath == "" {
				xcodeProjPath = path
			}
		} else {
			dsymPath = path
		}

		if xcodeProjPath != "" {
			if xcodeBuildOptions.Shared.ProjectRoot == "" {
				xcodeBuildOptions.Shared.ProjectRoot = ios.GetDefaultProjectRoot(xcodeProjPath, xcodeBuildOptions.Shared.ProjectRoot)
				logger.Info(fmt.Sprintf("Setting `--project-root` from Xcode project settings: %s", xcodeBuildOptions.Shared.ProjectRoot))
			}

			if xcodeBuildOptions.Shared.Scheme == "" {
				xcodeBuildOptions.Shared.Scheme, err = ios.GetDefaultScheme(xcodeProjPath)
				if err != nil {
					logger.Warn(fmt.Sprintf("Error determining default scheme: %s", err))
				}
			} else {
				_, err = ios.IsSchemeInPath(xcodeProjPath, xcodeBuildOptions.Shared.Scheme)
				if err != nil {
					logger.Warn(fmt.Sprintf("Scheme validation error: %s", err))
				}
			}

			if xcodeBuildOptions.Shared.Scheme != "" {
				buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjPath, xcodeBuildOptions.Shared.Scheme, xcodeBuildOptions.Shared.Scheme)
				if err != nil {
					logger.Warn(fmt.Sprintf("Error retrieving build settings: %s", err))
				}
			}

			if buildSettings != nil && dsymPath == "" {
				possibleDsymPath := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)
				if _, err = os.Stat(possibleDsymPath); err == nil {
					dsymPath = possibleDsymPath
					logger.Debug(fmt.Sprintf("Using dSYM path: %s", dsymPath))
				}
			}
		}

		if xcodeBuildOptions.Shared.ProjectRoot == "" {
			xcodeBuildOptions.Shared.ProjectRoot, _ = os.Getwd()
			logger.Info(fmt.Sprintf("Setting `--project-root` to current working directory: %s", xcodeBuildOptions.Shared.ProjectRoot))
		}

		if dsymPath == "" {
			return nil, "", nil, fmt.Errorf("No dSYM locations detected. Provide a valid dSYM path or Xcode project/workspace path")
		}

		dwarfInfo, tempDir, err = ios.FindDsymsInPath(dsymPath, xcodeBuildOptions.Shared.IgnoreEmptyDsym, xcodeBuildOptions.Shared.IgnoreMissingDwarf, logger)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return nil, "", tempDirs, fmt.Errorf("Error locating dSYM files: %w", err)
		}
		if len(dwarfInfo) == 0 {
			return nil, "", tempDirs, fmt.Errorf("No dSYM files found in: %s", dsymPath)
		}

		if plistPath == "" && opts.ApiKey == "" && buildSettings != nil {
			plistPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
		}
	}

	return dwarfInfo, plistPath, tempDirs, nil
}

func UploadDsyms(dsyms []*ios.DwarfInfo, plistPath, endpoint string, opts options.CLI, logger log.Logger) error {
	projectRoot := opts.Upload.XcodeBuild.Shared.ProjectRoot
	err := ios.ProcessDsymUpload(plistPath, endpoint, projectRoot, opts, dsyms, logger)
	if err != nil {
		return fmt.Errorf("Error uploading dSYM files: %w", err)
	}
	return nil
}
