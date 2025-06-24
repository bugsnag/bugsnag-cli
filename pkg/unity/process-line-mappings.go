package unity

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
)

func UploadUnityLineMappings(
	apiKey string,
	platform string,
	buildId string,
	applicationId string,
	appVersion string,
	platformVersion string,
	lineMappingFile string,
	projectRoot string,
	overwrite bool,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	uploadOptions := make(map[string]string)

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	if platform == "android" {
		if buildId != "" {
			uploadOptions["soBuildId"] = buildId
		}

		if platformVersion != "" {
			uploadOptions["appVersionCode"] = platformVersion
		}

	} else if platform == "ios" {
		if buildId != "" {
			uploadOptions["dsymUUID"] = buildId
		}

		if platformVersion != "" {
			uploadOptions["appBundleVersion"] = platformVersion
		}
	}

	if applicationId != "" {
		uploadOptions["appId"] = applicationId
	}

	if appVersion != "" {
		uploadOptions["appVersion"] = appVersion
	}

	if projectRoot != "" {
		uploadOptions["projectRoot"] = projectRoot
	}

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	fileFieldData := map[string]server.FileField{
		"mappingFile": server.LocalFile(lineMappingFile),
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
