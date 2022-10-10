package upload

import (
	"debug/elf"
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type DartSymbol struct {
	Path []string `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`
	AppVersion string `help:"(optional) the version of the application."`
	AppVersionCode string `help:"(optional) the version code for the application (Android only)."`
	AppBundleVersion string `help:"(optional) the bundle version for the application (iOS only)."`
	IosAppPath string `help:"(optional) the path to the built IOS app."`
}

func Dart(paths []string,
	appVersion string,
	appVersionCode string,
	appBundleVersion string,
	iosAppPath string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	apiKey string) error {

	// Build the file list from the path(s)
	log.Info("building file list...")

	fileList, err := utils.BuildFileList(paths)

	if err != nil {
		log.Error(" error building file list", 1)
	}

	log.Info("File list built...")

	for _, file := range fileList {

		// Process Android
		androidPlatform, _ := regexp.MatchString("android-([^;]*).symbols", file)

		if androidPlatform {
			log.Info("Processing android file " + file + "...")
			buildId, err := GetAndroidBuildId(file)

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
			continue
		}

		// Process IOS
		iosPlatform, _ := regexp.MatchString("ios-([^;]*).symbols", file)

		if iosPlatform {
			log.Info("Processing IOS file " + file + "...")

			// Builds the base path from the file that its processing
			regx := regexp.MustCompile(`/[^/]*$`)
			basePath := regx.ReplaceAllString(file, "")
			basePath = regx.ReplaceAllString(basePath, "")
			basePath = basePath + "/build/ios/iphoneos/"

			str, err := GetIosBuildId(basePath)

			if err != nil {
				log.Error(err.Error(), 1)
			}

			log.Info(str)

			continue
		}

		log.Warn("No files to process...")
	}

	return nil
}

// GetIosBuildId - Gets the build ID of an IOS App (Dwarf) file
func GetIosBuildId(basePath string) (string, error) {
	files, err := ioutil.ReadDir(basePath)

	if err != nil {
		return "", fmt.Errorf("unable to read files from dir " + err.Error())
	}

	for _,file := range files {
		if strings.Contains(file.Name(), ".app") && file.IsDir() {
			log.Info(basePath + file.Name() + "/Frameworks/App.framework/App")

			elfFile, err := elf.Open("/Users/josh.edney/repos/bugsnag-cli/examples/flutter/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App")

			if err != nil {
				return "", fmt.Errorf("unable to open elf file " + err.Error())
			}

			dwarfData, err := elfFile.DWARF()

			if err != nil {
				return "", fmt.Errorf("unable to get dwarf data " + err.Error())
			}

			entryReader := dwarfData.Reader()

			for {
				entry, err := entryReader.Next()

				if err != nil {
					return "", fmt.Errorf("unable to read entry from dwarf file " + err.Error())
				}

				fmt.Println(entry)
			}

		} else {
			return "", fmt.Errorf("unable to find IOS app")
		}
	}

	return "", nil
}

//GetAndroidBuildId - Gets the build ID of an Android symbol (elf) file
func GetAndroidBuildId(path string) (string, error){
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

	} else {
		return "", fmt.Errorf("no build id found")
	}

	return "", nil
}
