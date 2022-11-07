package upload

import (
	"debug/elf"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type DartSymbol struct {
	Path             utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`
	AppVersion       string            `help:"(optional) the version of the application."`
	AppVersionCode   string            `help:"(optional) the version code for the application (Android only)."`
	AppBundleVersion string            `help:"(optional) the bundle version for the application (iOS only)."`
	IosAppPath       string            `help:"(optional) the path to the built iOS app."`
}

func Dart(paths []string, appVersion string, appVersionCode string, appBundleVersion string, iosAppPath string, endpoint string, timeout int, retries int, overwrite bool, apiKey string) error {
	log.Info("Building file list from path")

	fileList, err := utils.BuildFileList(paths)

	if err != nil {
		log.Error("error building file list", 1)
	}

	log.Info("File list built")

	for _, file := range fileList {

		// Check if we're dealing with an android or iOS symbol file
		androidPlatform, _ := regexp.MatchString("android-([^;]*).symbols", file)
		iosPlatform, _ := regexp.MatchString("ios-([^;]*).symbols", file)

		// Start processing the android symbol file
		if androidPlatform {
			log.Info("Processing android symbol file: " + file)

			buildId, err := GetBuildIdFromElfFile(file)

			if err != nil {
				return err
			}

			// Build Upload options
			uploadOptions := BuildUploadOptions(apiKey, buildId, "android", overwrite, appVersion, appVersionCode)

			fileFieldName := "symbolFile"

			requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldName, file, timeout)

			if requestStatus != nil {
				return requestStatus
			}

			log.Success(file)

			continue
		}

		// Process iOS file
		if iosPlatform {
			log.Info("Processing iOS symbol file: " + file)

			if iosAppPath == "" {
				iosAppPath, err = GetIosAppPath(file)

				if err != nil {
					return err
				}
			}

			arch, err := GetArchFromElfFile(file)

			if err != nil {
				return err
			}

			buildId, err := DwarfDumpUuid(file, iosAppPath, arch)

			if err != nil {
				return err
			}

			// Build Upload options
			uploadOptions := BuildUploadOptions(apiKey, buildId, "ios", overwrite, appVersion, appBundleVersion)

			fileFieldName := "symbolFile"

			requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldName, file, timeout)

			if requestStatus != nil {
				return requestStatus
			}

			log.Success(file)

			continue
		}
		log.Info("Skipping " + file)

	}

	return nil
}

// BuildUploadOptions - Builds the upload options for processing dart files
func BuildUploadOptions(apiKey string, uuid string, platform string, overwrite bool, appVersion string, appExtraVersion string) map[string]string {
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = apiKey

	uploadOptions["buildId"] = uuid

	uploadOptions["platform"] = platform

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	if platform == "ios" {
		if appVersion != "" {
			uploadOptions["appVersion"] = appVersion
		}

		if appExtraVersion != "" {
			uploadOptions["AppBundleVersion"] = appExtraVersion
		}
	}

	if platform == "android" {
		if appVersion != "" {
			uploadOptions["appVersion"] = appVersion
		}

		if appExtraVersion != "" {
			uploadOptions["appVersionCode"] = appExtraVersion
		}
	}

	return uploadOptions
}

// ReadElfFile - Gets all data from the symbol file
func ReadElfFile(symbolFile string) (*elf.File, error) {
	file, err := os.OpenFile(symbolFile, os.O_RDONLY, 0)

	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	_elf, err := elf.NewFile(file)

	if err != nil {
		return nil, fmt.Errorf("error reading symbol file: %w", err)
	}

	return _elf, nil
}

// GetBuildIdFromElfFile - Gets the build ID from the symbol file
func GetBuildIdFromElfFile(symbolFile string) (string, error) {
	elfData, err := ReadElfFile(symbolFile)

	if err != nil {
		return "", err
	}

	if sect := elfData.Section(".note.gnu.build-id"); sect != nil {
		data, err := sect.Data()

		if err != nil {
			return "", fmt.Errorf("error reading symbol file")
		}

		buildId := fmt.Sprintf("%x", data[16:])

		return buildId, nil
	}

	return "", fmt.Errorf("unable to find buildId")
}

// GetArchFromElfFile - Gets the Arch from the symbol file to help with getting the UUID from the built iOS app
func GetArchFromElfFile(symbolFile string) (string, error) {
	elfData, err := ReadElfFile(symbolFile)

	if err != nil {
		return "", err
	}

	var arch string

	switch elfData.Machine {
	case EM_AARCH64:
		arch = "arm64"
	case EM_386:
		arch = "x86"
	case EM_X86_64:
		arch = "x86_64"
	case EM_ARM:
		arch = "armv7"
	}

	return arch, nil
}

// GetIosAppPath - Gets the path to the built iOS app relative to the symbol files
func GetIosAppPath(symbolFile string) (string, error) {
	sampleRegexp := regexp.MustCompile(`/[^/]*/[^/]*$`)
	basePath := filepath.Join(sampleRegexp.ReplaceAllString(symbolFile, "") + "/build/ios/iphoneos")
	files, err := ioutil.ReadDir(basePath)

	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".app") && file.IsDir() {
			iosAppPath := filepath.Join(basePath + "/" + file.Name() + "/Frameworks/App.framework/App")
			return iosAppPath, nil
		}
	}

	return "", fmt.Errorf("unable to find iOS app path, try adding --ios-app-path")
}

// DwarfDumpUuid - Gets the UUID/Build ID from the Dwarf debug info of a file for a given Arch
func DwarfDumpUuid(symbolFile string, dwarfFile string, arch string) (string, error) {
	dwarfDumpLocation, err := exec.LookPath("dwarfdump")
	uuidArray := make(map[string]string)

	if err != nil {
		return "", fmt.Errorf("unable to find dwarfdump on system: %w", err)
	}

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
