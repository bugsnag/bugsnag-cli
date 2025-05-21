package upload

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/elf"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

var androidSymbolFileRegex = regexp.MustCompile("android-([^;]*).symbols")
var iosSymbolFileRegex = regexp.MustCompile("ios-([^;]*).symbols")

func Dart(options options.CLI, endpoint string, logger log.Logger) error {
	dartOptions := options.Upload.DartSymbol
	fileList, err := utils.BuildFileList(dartOptions.Path)

	if err != nil {
		logger.Fatal("error building file list")
	}

	for _, file := range fileList {

		// Check if we're dealing with an android or iOS symbol file
		isAndroidPlatform := androidSymbolFileRegex.MatchString(file)
		isIosPlatform := iosSymbolFileRegex.MatchString(file)

		// Start processing the android symbol file
		if isAndroidPlatform {
			logger.Info(fmt.Sprintf("Processing android symbol file: %s", file))

			var buildId string
			buildId, err = elf.GetBuildId(file)
			if err != nil {
				return err
			}

			// Build Upload options
			uploadOptions := utils.BuildDartUploadOptions(options.ApiKey, buildId, "android", options.Upload.Overwrite, dartOptions.VersionName, dartOptions.VersionCode)

			fileFieldData := make(map[string]server.FileField)
			fileFieldData["symbolFile"] = server.LocalFile(file)

			err := server.ProcessFileRequest(endpoint+"/dart-symbol", uploadOptions, fileFieldData, file, options, logger)

			if err != nil {

				return err
			}

			continue
		}

		// Process iOS file
		if isIosPlatform {
			logger.Info(fmt.Sprintf("Processing iOS symbol file: %s", file))

			iosAppPath := string(dartOptions.IosAppPath)
			if iosAppPath == "" {
				iosAppPath, err = GetIosAppPath(file)

				if err != nil {
					return err
				}
			}

			var arch string
			arch, err = elf.GetArch(file)
			if err != nil {
				return err
			}

			var buildId string
			buildId, err = DwarfDumpUuid(file, iosAppPath, arch)
			if err != nil {
				return err
			}

			// Build Upload options
			uploadOptions := utils.BuildDartUploadOptions(options.ApiKey, buildId, "ios", options.Upload.Overwrite, dartOptions.VersionName, dartOptions.BundleVersion)

			fileFieldData := make(map[string]server.FileField)
			fileFieldData["symbolFile"] = server.LocalFile(file)

			if options.DryRun {
				err = nil
			} else {
				err = server.ProcessFileRequest(endpoint+"/dart-symbol", uploadOptions, fileFieldData, file, options, logger)
			}

			if err != nil {

				return err
			}

			continue
		}
		logger.Debug(fmt.Sprintf("Skipping %s - unsupported platform", file))
	}

	return nil
}

// GetIosAppPath - Gets the path to the built iOS app relative to the symbol files
func GetIosAppPath(symbolFile string) (string, error) {
	sampleRegexp := regexp.MustCompile(`/[^/]*/[^/]*$`)
	basePath := filepath.Join(sampleRegexp.ReplaceAllString(symbolFile, "") + "/build/ios/iphoneos")
	files, err := os.ReadDir(basePath)

	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".app") && file.IsDir() {
			iosAppPath := filepath.Join(basePath + "/" + file.Name() + "/Frameworks/App.framework/App")

			_, err = os.Stat(iosAppPath)

			if errors.Is(err, os.ErrNotExist) {
				return "", err
			}

			return iosAppPath, nil
		}
	}

	return "", fmt.Errorf("unable to find iOS app path, try adding --ios-app-path")
}

// DwarfDumpUuid - Gets the UUID/Build ID from the Dwarf debug info of a file for a given Arch
func DwarfDumpUuid(symbolFile string, dwarfFile string, arch string) (string, error) {
	dwarfDumpLocation, err := exec.LookPath("dwarfdump")

	if err != nil {
		dwarfDumpLocation, err = exec.LookPath("llvm-dwarfdump")

		if err != nil {
			return "", fmt.Errorf("unable to find dwarfdump on system: %w", err)
		}
	}

	uuidArray := make(map[string]string)
	cmd := exec.Command(dwarfDumpLocation, "--uuid", dwarfFile, "--arch", arch)
	output, _ := cmd.CombinedOutput()
	outputArray := strings.Fields(string(output))

	uuidArray[outputArray[2]] = outputArray[1]

	for key, value := range uuidArray {
		uuidArch := strings.Replace(key, "(", "", -1)
		uuidArch = strings.Replace(uuidArch, ")", "", -1)

		if strings.Contains(symbolFile, uuidArch) {
			return value, nil
		}
	}

	return "", fmt.Errorf("unable to find matching UUID")
}
