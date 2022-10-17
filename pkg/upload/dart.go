package upload

import (
	"debug/elf"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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
	IosAppPath       string            `help:"(optional) the path to the built IOS app."`
}

func Dart(paths []string, appVersion string, appVersionCode string, appBundleVersion string, iosAppPath string, endpoint string, timeout int, retries int, overwrite bool, apiKey string) error {
	log.Info("Building file list from path")

	fileList, err := utils.BuildFileList(paths)

	if err != nil {
		log.Error("error building file list", 1)
	}

	log.Info("File list built")

	for _, file := range fileList {

		// Check if we're dealing with an android or IOS symbol file
		androidPlatform, _ := regexp.MatchString("android-([^;]*).symbols", file)
		iosPlatform, _ := regexp.MatchString("ios-([^;]*).symbols", file)

		// Start processing the android symbol file
		if androidPlatform {
			log.Info("Processing android symbol file: " + file)

			buildId, err := ReadElfBuildId(file)

			if err != nil {
				return err
			}

			// Build Upload options
			uploadOptions := make(map[string]string)

			uploadOptions["apiKey"] = apiKey

			uploadOptions["buildId"] = buildId

			uploadOptions["platform"] = "android"

			if overwrite {
				uploadOptions["overwrite"] = "true"
			}

			if appVersion != "" {
				uploadOptions["appVersion"] = appVersion
			}

			if appVersionCode != "" {
				uploadOptions["appVersionCode"] = appVersionCode
			}

			fileFieldName := "symbolFile"

			req, err := server.BuildFileRequest(endpoint, uploadOptions, fileFieldName, file)

			if err != nil {
				return fmt.Errorf("error building file request: %w", err)
			}

			res, err := server.SendRequest(req, timeout)

			if err != nil {
				return fmt.Errorf("error sending file request: %w", err)
			}

			b, err := io.ReadAll(res.Body)

			if err != nil {
				return fmt.Errorf("error reading body from response: %w", err)
			}

			if res.Status != "202 Accepted" {
				return fmt.Errorf("%s : %s", res.Status, string(b))
			}

			log.Success(file)

			continue
		}

		// Process IOS file
		if iosPlatform {
			log.Info("Processing IOS symbol file: " + file)

			if iosAppPath == "" {
				iosAppPath, err = GetIosAppPath(file)

				if err != nil {
					return err
				}
			}

			uuid, err := DwarfDumpUuid(file, iosAppPath)

			if err != nil {
				return err
			}

			// Build Upload options
			uploadOptions := make(map[string]string)

			uploadOptions["apiKey"] = apiKey

			uploadOptions["buildId"] = uuid

			uploadOptions["platform"] = "ios"

			if overwrite {
				uploadOptions["overwrite"] = "true"
			}

			if appVersion != "" {
				uploadOptions["appVersion"] = appVersion
			}

			if appBundleVersion != "" {
				uploadOptions["AppBundleVersion"] = appBundleVersion
			}

			fileFieldName := "symbolFile"

			req, err := server.BuildFileRequest(endpoint, uploadOptions, fileFieldName, file)

			if err != nil {
				return fmt.Errorf("error building file request: %w", err)
			}

			res, err := server.SendRequest(req, timeout)

			if err != nil {
				return fmt.Errorf("error sending file request: %w", err)
			}

			b, err := io.ReadAll(res.Body)

			if err != nil {
				return fmt.Errorf("error reading body from response: %w", err)
			}

			if res.Status != "202 Accepted" {
				return fmt.Errorf("%s : %s", res.Status, string(b))
			}

			log.Success(file)

			continue
		}
		log.Info("Skipping " + file)

	}

	return nil
}

// GetIosAppPath - Gets the base path to the built IOS app
func GetIosAppPath(symbolFile string) (string, error) {
	sampleRegexp := regexp.MustCompile(`/[^/]*/[^/]*$`)
	basePath := sampleRegexp.ReplaceAllString(symbolFile, "") + "/build/ios/iphoneos/"

	files, err := ioutil.ReadDir(basePath)

	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".app") && file.IsDir() {
			iosAppPath := basePath + file.Name() + "/Frameworks/App.framework/App"
			return iosAppPath, nil
		}
	}

	return "", fmt.Errorf("unable to find IOS app path, try adding --ios-app-path")
}

// DwarfDumpUuid - Gets the UUID from the Dwarf debug info of a file
func DwarfDumpUuid(symbolFile string, dwarfFile string) (string, error) {
	dwarfDumpLocation, err := exec.LookPath("dwarfdump")
	uuidArray := make(map[string]string)

	if err != nil {
		return "", fmt.Errorf("unable to find dwarfdump on system: %w", err)
	}

	cmd := exec.Command(dwarfDumpLocation, "--uuid", dwarfFile)
	output, _ := cmd.CombinedOutput()
	outputArray := strings.Fields(string(output))

	uuidArray[outputArray[2]] = outputArray[1]
	uuidArray[outputArray[6]] = outputArray[5]

	for key, value := range uuidArray {
		uuidArch := strings.Replace(key, "(", "", -1)
		uuidArch = strings.Replace(uuidArch, ")", "", -1)

		if strings.Contains(symbolFile, uuidArch) {
			return value, nil
		}
	}

	return "", fmt.Errorf("unable to find matching UUID")
}

func ReadElfBuildId(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)

	if err != nil {
		return "", fmt.Errorf("unable to open file")
	}

	elfData, err := elf.NewFile(file)

	if err != nil {
		return "", fmt.Errorf("error reading symbol file")
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
