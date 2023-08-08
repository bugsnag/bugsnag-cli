package upload

import (
	"debug/elf"
	"errors"
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
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	IosAppPath    string            `help:"(optional) the path to the built iOS app."`
	VersionName   string            `help:"The version of the application." xor:"app-version,version-name"`
	VersionCode   string            `help:"The version code for the application (Android only)." xor:"app-version-code,version-code"`
	BundleVersion string            `help:"The bundle version for the application (iOS only)." xor:"app-bundle-version,bundle-version"`
}

func Dart(paths []string, version string, versionCode string, bundleVersion string, iosAppPath string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool, dryRun bool) error {

	log.Info("Building file list from path")

	fileList, err := utils.BuildFileList(paths)
	numberOfFiles := len(fileList)

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
			uploadOptions := utils.BuildDartUploadOptions(apiKey, buildId, "android", overwrite, version, versionCode)

			fileFieldData := make(map[string]string)
			fileFieldData["symbolFile"] = file

			requestStatus := server.ProcessRequest(endpoint+"/dart-symbol", uploadOptions, fileFieldData, timeout, file, dryRun)

			if requestStatus != nil {
				if numberOfFiles > 1 && failOnUploadError {
					return requestStatus
				} else {
					log.Warn(requestStatus.Error())
				}
			} else {
				log.Success(file)
			}

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
			uploadOptions := utils.BuildDartUploadOptions(apiKey, buildId, "ios", overwrite, version, bundleVersion)

			fileFieldData := make(map[string]string)
			fileFieldData["symbolFile"] = file

			if dryRun {
				err = nil
			} else {
				err = server.ProcessRequest(endpoint+"/dart-symbol", uploadOptions, fileFieldData, timeout, file, dryRun)
			}

			if err != nil {
				if numberOfFiles > 1 && failOnUploadError {
					return err
				} else {
					log.Warn(err.Error())
				}
			} else {
				log.Success(file)
			}

			continue
		}
		log.Info("Skipping " + file)
	}

	return nil
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
	case elf.EM_AARCH64:
		arch = "arm64"
	case elf.EM_386:
		arch = "x86"
	case elf.EM_X86_64:
		arch = "x86_64"
	case elf.EM_ARM:
		arch = "armv7"
	}

	if arch == "" {
		return "", fmt.Errorf("unable to find arch type")
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

			_, err := os.Stat(iosAppPath)

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
