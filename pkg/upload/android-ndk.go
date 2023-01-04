package upload

import (
	"fmt"
	"os"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidNdkMapping struct {
	AndroidNdkRoot  string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifestPath string            `help:"(required) Path to directory or file to upload" type:"path"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot     string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}

//soFile - the path to the shared object file.
//apiKey - your Bugsnag integration API key for this application.
//appId (optional) - the Android applicationId for this application.
//versionCode (optional) - the Android versionCode for this application release.
//sharedObjectName (optional) - the name of the shared object that the symbols are for.
//versionName (optional) - the Android versionName for this application release.
//projectRoot (optional) - a path to remove from the beginning of the filenames in the mapping file
//overwrite (optional) - overwrite any existing mappings for this version of your app.

func ProcessAndroidNDK(paths []string, androidNdkRoot string, appManifestPath string, configuration string, projectRoot string, buildUuid string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {

	// Check to see if we've set the AppManifestPath
	if appManifestPath == "" {

		appManifestPath = ../../app/build/intermediates/merged-manifests/[variant]/AndroidManifest.xml
	}

	if configuration == "" {

	}

	// Check to see if we've got the project root set.
	if projectRoot == "" {
		return fmt.Errorf("missing `--project-root` option from the command line")
	}

	log.Info("Building file list from path")
	fileList, err := utils.BuildFileList(paths)
	numberOfFiles := len(fileList)

	fmt.Println(fileList)
	fmt.Println(string(numberOfFiles))

	if err != nil {
		log.Error("error building file list", 1)
	}

	log.Info("File list built")

	// Check if we have Android NDK root set
	log.Info("Setting Android NDK Root")

	ndkPath, err := GetNDKRootPath(androidNdkRoot)

	if err != nil {
		return err
	}

	log.Info("Using NDK path: " + ndkPath)

	fmt.Println(appManifestPath)
	//
	//log.Info("Locating ObjCopy within NDK path")
	//
	//objcopyPath, err := BuildObjCopyPath(ndkPath)
	//
	//if err != nil {
	//	return err
	//}
	//
	//log.Info("Found ObjCopy in NDK path")
	//
	//for _, file := range fileList {
	//	if filepath.Ext(file) == ".so" {
	//		log.Info("Converting " + filepath.Base(file) + " using objcopy")
	//		outputFile, err := ObjCopy(objcopyPath, file)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//}
	return nil
}

//	uploadFileOptions := make(map[string]string)
//
//	for _, file := range fileList {
//		if filepath.Ext(file) == ".so" {
//			log.Info("Converting " + filepath.Base(file) + " using objcopy")
//			outputFile, err := ObjCopy(objcopyPath, file)
//			if err != nil {
//				return err
//			}
//			uploadFileOptions[filepath.Base(file)] = outputFile
//		}
//	}
//
//	for key, value := range uploadFileOptions {
//		uploadOptions := BuildNDKUploadOptions(apiKey, buildId, "android", overwrite, appVersion, appVersionCode)
//		//			requestStatus := server.ProcessRequest(endpoint, uploadOptions, "symbolFile", file, timeout)
//		log.Info(key + ":" + value)
//	}
//
//	return nil
//}

// GetNDKRootPath - Checks if we're passing the NDK root path via the options.
// Defaults to OS environment variable if not passed via the options.
func GetNDKRootPath(ndkPath string) (string, error) {
	if ndkPath == "" {
		androidNdkRootValue, androidNdkRootPresent := os.LookupEnv("ANDROID_NDK_ROOT")
		if androidNdkRootPresent {
			ndkPath = androidNdkRootValue
		} else {
			return "", fmt.Errorf("unable to find ANDROID_NDK_ROOT")
		}
	}

	if _, err := os.Stat(ndkPath); os.IsNotExist(err) {
		return ndkPath, fmt.Errorf("ANDROID_NDK_ROOT path does not exist on the system: " + ndkPath)
	}
	return ndkPath, nil
}

// BuildObjCopyPath - Builds the path to the ObjCopy binary within the NDK root path
//func BuildObjCopyPath(ndkPath string) (string, error) {
//	ndkVersion, err := GetNdkVersion(ndkPath)
//	fmt.Println(ndkVersion)
//	if err != nil {
//		return "", fmt.Errorf("unable to determine ndk version from path")
//	}
//
//	if ndkVersion < 24 {
//		directoryPattern := filepath.Join(ndkPath, "/toolchains/x86_64-4.9/prebuilt/*/bin")
//		directoryMatches, err := filepath.Glob(directoryPattern)
//		if err != nil {
//			return "", err
//		}
//		if directoryMatches == nil {
//			return "", fmt.Errorf("Unable to find objcopy within ANDROID_NDK_ROOT: " + ndkPath)
//		}
//
//		if runtime.GOOS == "windows" {
//			return filepath.Join(directoryMatches[0], "x86_64-linux-android-objcopy.exe"), nil
//		}
//
//		return filepath.Join(directoryMatches[0], "x86_64-linux-android-objcopy"), nil
//	} else {
//		directoryPattern := filepath.Join(ndkPath, "/toolchains/llvm/prebuilt/*/bin")
//		directoryMatches, err := filepath.Glob(directoryPattern)
//		if err != nil {
//			return "", err
//		}
//
//		if directoryMatches == nil {
//			return "", fmt.Errorf("Unable to find objcopy within ANDROID_NDK_ROOT: " + ndkPath)
//		}
//
//		if runtime.GOOS == "windows" {
//			return filepath.Join(directoryMatches[0], "llvm-objcopy.exe"), nil
//		}
//
//		return filepath.Join(directoryMatches[0], "llvm-objcopy"), nil
//	}
//	return "", nil
//}
//
//// ObjCopy - Uses ObjCopy
//func ObjCopy(objcopyPath string, file string) (string, error) {
//
//	objcopyLocation, err := exec.LookPath(objcopyPath)
//
//	if err != nil {
//		return "", err
//	}
//
//	cmd := exec.Command(objcopyLocation, "--compress-debug-sections=zlib", "--only-keep-debug", file, file+"-out")
//
//	_, err = cmd.CombinedOutput()
//
//	return file + "-out", nil
//}
//
//func GetNdkVersion(ndkPath string) (int, error) {
//	ndkVersion := strings.Split(filepath.Base(ndkPath), ".")
//	ndkIntVersion, err := strconv.Atoi(ndkVersion[0])
//	if err != nil {
//		return 0, err
//	}
//	return ndkIntVersion, nil
//}
//
//// BuildNDKUploadOptions BuildUploadOptions - Builds the upload options for processing NDK files
//func BuildNDKUploadOptions(apiKey string, uuid string, platform string, overwrite bool, appVersion string, appExtraVersion string) map[string]string {
//	//-F versionCode=123 \
//	//-F appId=com.example.android.app \
//	//-F sharedObjectName=libmy-ndk-library.so \
//	//-F versionName=2.3.0
//
//	uploadOptions := make(map[string]string)
//
//	uploadOptions["apiKey"] = apiKey
//
//	uploadOptions["versionCode"] = uuid
//
//	uploadOptions["versionName"] = uuid
//
//	uploadOptions["appId"] = platform
//
//	if overwrite {
//		uploadOptions["overwrite"] = "true"
//	}
//
//	return uploadOptions
//}


