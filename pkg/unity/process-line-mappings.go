package unity

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func UploadAndroidLineMappings(
	lineMappingFile string,
	soBuildId string,
	endpoint string,
	options options.CLI,
	manifestData map[string]string,
	logger log.Logger,
) error {
	opts := utils.UnityLineMappingOptions{
		APIKey:         manifestData["apiKey"],
		AppID:          manifestData["applicationId"],
		AppVersion:     manifestData["versionName"],
		AppVersionCode: manifestData["versionCode"],
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
