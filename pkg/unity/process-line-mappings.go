package unity

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// UnityLineMappingOptions holds shared upload metadata for Unity line mapping uploads.
// Includes platform-specific fields for Android and iOS.
type UnityLineMappingOptions struct {
	APIKey           string
	AppID            string
	AppVersion       string
	AppVersionCode   string // Android only
	AppBundleVersion string // iOS only
	SOBuildID        string // Android only
	DSYMUUUID        string // iOS only
	ProjectRoot      string
	Overwrite        bool
}

// UploadAndroidLineMappings uploads an Android Unity line mapping file to the Bugsnag server.
// It builds the upload options and submits the file using the provided endpoint and logger.
func UploadAndroidLineMappings(
	lineMappingFile string,
	soBuildId string,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	opts := UnityLineMappingOptions{
		APIKey:         options.ApiKey,
		AppID:          options.Upload.UnityAndroid.ApplicationId,
		AppVersion:     options.Upload.UnityAndroid.VersionName,
		AppVersionCode: options.Upload.UnityAndroid.VersionCode,
		SOBuildID:      soBuildId,
		ProjectRoot:    options.Upload.UnityAndroid.ProjectRoot,
		Overwrite:      options.Upload.UnityAndroid.Overwrite,
	}

	fileFieldData := map[string]server.FileField{
		"mappingFile": server.LocalFile(lineMappingFile),
	}

	uploadOptions, err := utils.BuildUnityLineMappingUploadOptions(opts)
	if err != nil {
		return err
	}

	return server.ProcessFileRequest(
		endpoint+"/unity-line-mappings",
		uploadOptions,
		fileFieldData,
		lineMappingFile,
		options,
		logger,
	)
}

// UploadIosLineMappings uploads an iOS Unity line mapping file to the Bugsnag server.
// It builds the upload options and submits the file using the provided endpoint and logger.
func UploadIosLineMappings(
	lineMappingFile string,
	dsymUuid string,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	opts := UnityLineMappingOptions{
		APIKey:           options.ApiKey,
		AppID:            options.Upload.UnityIos.ApplicationId,
		AppVersion:       options.Upload.UnityIos.VersionName,
		AppBundleVersion: options.Upload.UnityIos.BundleVersion,
		DSYMUUUID:        dsymUuid,
		ProjectRoot:      options.Upload.UnityIos.Shared.ProjectRoot,
		Overwrite:        options.Upload.UnityIos.Overwrite,
	}

	fileFieldData := map[string]server.FileField{
		"mappingFile": server.LocalFile(lineMappingFile),
	}

	uploadOptions, err := utils.BuildUnityLineMappingUploadOptions(opts)
	if err != nil {
		return err
	}

	return server.ProcessFileRequest(
		endpoint+"/unity-line-mappings",
		uploadOptions,
		fileFieldData,
		lineMappingFile,
		options,
		logger,
	)
}
