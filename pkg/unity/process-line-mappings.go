package unity

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)
import "github.com/bugsnag/bugsnag-cli/pkg/options"

func UploadAndroidLineMappings(
	lineMappingFile string,
	apiKey string,
	soBuildId string,
	applicationId string,
	versionName string,
	versionCode string,
	projectRoot string,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	fileFieldData := make(map[string]server.FileField)

	uploadOptions, err := utils.BuildUnityAndroidLineMappingUploadOptions(apiKey, soBuildId, applicationId, versionName, versionCode, projectRoot, options.Upload.UnityAndroid.Overwrite)

	if err != nil {
		return err
	}

	fileFieldData["mappingFile"] = server.LocalFile(lineMappingFile)

	err = server.ProcessFileRequest(endpoint+"/unity-line-mappings", uploadOptions, fileFieldData, lineMappingFile, options, logger)

	if err != nil {
		return err
	}

	return nil
}
