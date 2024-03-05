package upload

import (
	"errors"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type Dsym struct {
	Scheme             string      `help:"The name of the scheme to use when building the application."`
	XcodeProject       utils.Path  `help:"Path to the dSYM" type:"path"`
	Plist              utils.Path  `help:"Path to the Info.plist file" type:"path"`
	IgnoreMissingDwarf bool        `help:"Throw warnings instead of errors when a dSYM with missing DWARF data is found"`
	IgnoreEmptyDsym    bool        `help:"Throw warnings instead of errors when a *.dSYM file is found, rather than the expected *.dSYM directory"`
	Path               utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessDsym(apiKey string, xcodeProjectPath string, scheme string, plistPath string, paths []string, endpoint string, overwrite bool, timeout int, retries int, dryRun bool, failOnUpload bool) error {
	var buildSettings *ios.XcodeBuildSettings
	var dsymPath string
	var plistData *ios.PlistData
	var uploadOptions map[string]string

	var dwarfInfo []*ios.DwarfInfo
	var tempDirs []string
	var err error

	for _, path := range paths {
		if utils.IsDir(path) {
			if filepath.Ext(path) == ".dSYM" {

				if apiKey == "" {
					if xcodeProjectPath != "" {
						// If scheme is set explicitly, check if it exists
						if scheme != "" {
							_, err := ios.IsSchemeInPath(xcodeProjectPath, scheme)
							if err != nil {
								log.Warn(err.Error())
							}

						} else {
							// Otherwise, try to find it
							var err error
							scheme, err = ios.GetDefaultScheme(xcodeProjectPath)
							if err != nil {
								log.Warn(err.Error())
							}
						}

						if scheme != "" {
							var err error
							buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjectPath, scheme)
							if err != nil {
								return err
							}
						}

						// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
						if plistPath == "" && apiKey == "" {
							if buildSettings != nil {
								plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
								if utils.FileExists(plistPathExpected) {
									plistPath = plistPathExpected
									log.Info("Found Info.plist at expected location: " + plistPath)
								} else {
									log.Info("No Info.plist found at expected location: " + plistPathExpected)
								}
							}
						}

						// If the Info.plist path is defined and we still don't know the apiKey or verionName, try to extract them from it
						if plistPath != "" && apiKey == "" {
							// Read data from the plist
							var err error
							plistData, err = ios.GetPlistData(plistPath)
							if err != nil {
								return err
							}

							if apiKey == "" {
								apiKey = plistData.BugsnagProjectDetails.ApiKey
								if apiKey != "" {
									log.Info("Using API key from Info.plist: " + apiKey)
								}
							}
						}
					} else {
						return errors.New("missing api key, please specify using `--api-key`")
					}
				}

				log.Info("Uploading dSYM: " + path)

				uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, overwrite)
				if err != nil {
					return err
				}

				fileFieldData := make(map[string]string)
				fileFieldData["dsym"] = filepath.Join(path, "Contents", "Resources", "DWARF", utils.GetBaseWithoutExt(filepath.Base(path)))

				err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, path, dryRun)

				if err != nil {
					if strings.Contains(err.Error(), "404 Not Found") {
						err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, path, dryRun)
					}
				}

				if err != nil {
					if failOnUpload {
						return err
					} else {
						log.Warn(err.Error())
					}
				} else {
					log.Success("Uploaded dSYM: " + path)
				}

			} else {
				log.Info("Processing project directory: " + path)
				if xcodeProjectPath == "" {
					xcodeProjectPath, _ = ios.FindProjectOrWorkspaceInPath(path)
				}

				// If scheme is set explicitly, check if it exists
				if scheme != "" {
					_, err := ios.IsSchemeInPath(xcodeProjectPath, scheme)
					if err != nil {
						log.Warn(err.Error())
					}

				} else {
					// Otherwise, try to find it
					scheme, err = ios.GetDefaultScheme(xcodeProjectPath)
					if err != nil {
						log.Warn(err.Error())
					}

				}

				if scheme != "" {
					buildSettings, err = ios.GetXcodeBuildSettings(path, scheme)
					if err != nil {
						return err
					}
				}

				if buildSettings != nil {
					dsymPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)

					_, err := os.Stat(dsymPath)
					if err == nil {
						log.Info("Using dSYM path: " + dsymPath)
					}
				}

				if dsymPath != "" {
					var tempDir string
					dwarfInfo, tempDir, _ = ios.FindDsymsInPath(dsymPath, false, false)
					tempDirs = append(tempDirs, tempDir)
				}
				if len(dwarfInfo) == 0 {
					return errors.New("No dSYM files found")
				}

				// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
				if plistPath == "" && apiKey == "" {
					if buildSettings != nil {
						plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
						if utils.FileExists(plistPathExpected) {
							plistPath = plistPathExpected
							log.Info("Found Info.plist at expected location: " + plistPath)
						} else {
							log.Info("No Info.plist found at expected location: " + plistPathExpected)
						}
					}
				}

				// If the Info.plist path is defined and we still don't know the apiKey or verionName, try to extract them from it
				if plistPath != "" && apiKey == "" {
					// Read data from the plist
					plistData, err = ios.GetPlistData(plistPath)
					if err != nil {
						return err
					}

					if apiKey == "" {
						apiKey = plistData.BugsnagProjectDetails.ApiKey
						if apiKey != "" {
							log.Info("Using API key from Info.plist: " + apiKey)
						}
					}
				}

				for _, dsym := range dwarfInfo {
					dsymInfo := "(UUID: " + dsym.UUID + ", Name: " + dsym.Name + ", Arch: " + dsym.Arch + ")"
					log.Info("Uploading dSYM " + dsymInfo)

					uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, overwrite)
					if err != nil {
						return err
					}

					fileFieldData := make(map[string]string)
					fileFieldData["dsym"] = filepath.Join(dsym.Location, dsym.Name)

					err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)

					if err != nil {
						if strings.Contains(err.Error(), "404 Not Found") {
							err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)
						}
					}

					if err != nil {
						if failOnUpload {
							return err
						} else {
							log.Warn(err.Error())
						}
					} else {
						log.Success("Uploaded dSYM: " + dsym.Name)
					}
				}
			}
		} else if filepath.Ext(path) == ".zip" {
			log.Info("Processing zip file: " + path)

			// Unzip the file
			var tempDir string
			dwarfInfo, tempDir, _ = ios.FindDsymsInPath(path, false, false)
			tempDirs = append(tempDirs, tempDir)

			if apiKey == "" {
				if xcodeProjectPath != "" {
					// If scheme is set explicitly, check if it exists
					if scheme != "" {
						_, err := ios.IsSchemeInPath(xcodeProjectPath, scheme)
						if err != nil {
							log.Warn(err.Error())
						}

					} else {
						// Otherwise, try to find it
						var err error
						scheme, err = ios.GetDefaultScheme(xcodeProjectPath)
						if err != nil {
							log.Warn(err.Error())
						}
					}

					if scheme != "" {
						var err error
						buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjectPath, scheme)
						if err != nil {
							return err
						}
					}

					// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
					if plistPath == "" && apiKey == "" {
						if buildSettings != nil {
							plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
							if utils.FileExists(plistPathExpected) {
								plistPath = plistPathExpected
								log.Info("Found Info.plist at expected location: " + plistPath)
							} else {
								log.Info("No Info.plist found at expected location: " + plistPathExpected)
							}
						}
					}

					// If the Info.plist path is defined and we still don't know the apiKey or verionName, try to extract them from it
					if plistPath != "" && apiKey == "" {
						// Read data from the plist
						var err error
						plistData, err = ios.GetPlistData(plistPath)
						if err != nil {
							return err
						}

						if apiKey == "" {
							apiKey = plistData.BugsnagProjectDetails.ApiKey
							if apiKey != "" {
								log.Info("Using API key from Info.plist: " + apiKey)
							}
						}
					}
				} else {
					return errors.New("missing api key, please specify using `--api-key`")
				}
			}

			for _, dsym := range dwarfInfo {
				dsymInfo := "(UUID: " + dsym.UUID + ", Name: " + dsym.Name + ", Arch: " + dsym.Arch + ")"
				log.Info("Uploading dSYM " + dsymInfo)

				uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, overwrite)
				if err != nil {
					return err
				}

				fileFieldData := make(map[string]string)
				fileFieldData["dsym"] = filepath.Join(dsym.Location, dsym.Name)

				err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)

				if err != nil {
					if strings.Contains(err.Error(), "404 Not Found") {
						err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)
					}
				}

				if err != nil {
					if failOnUpload {
						return err
					} else {
						log.Warn(err.Error())
					}
				} else {
					log.Success("Uploaded dSYM: " + dsym.Name)
				}
			}

		} else {
			return errors.New("Invalid file type, must be a .dSYM directory or a .zip file")
		}
	}

	return nil
}
