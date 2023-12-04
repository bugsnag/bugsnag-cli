package build

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type SourceControl struct {
	Provider   string
	Repository string
	Revision   string
}

type CreateBuildInfo struct {
	ApiKey            string            `json:"apiKey,omitempty"`
	AppVersionCode    string            `json:"appVersionCode,omitempty"`
	AppBundleVersion  string            `json:"appBundleVersion,omitempty"`
	SourceControl     SourceControl     `json:"sourceControl,omitempty"`
	BuilderName       string            `json:"builderName,omitempty"`
	ReleaseStage      string            `json:"releaseStage,omitempty"`
	AppVersion        string            `json:"appVersion,omitempty"`
	AutoAssignRelease *bool             `json:"autoAssignRelease,omitempty"`
	MetaData          map[string]string `json:"metadata,omitempty"`
}

func (opts CreateBuildInfo) Override(base CreateBuildInfo) CreateBuildInfo {
	return CreateBuildInfo{
		ApiKey:           utils.ThisOrThat(opts.ApiKey, base.ApiKey).(string),
		AppVersionCode:   utils.ThisOrThat(opts.AppVersionCode, base.AppVersionCode).(string),
		AppBundleVersion: utils.ThisOrThat(opts.AppBundleVersion, base.AppBundleVersion).(string),
		SourceControl: SourceControl{
			Provider:   utils.ThisOrThat(opts.SourceControl.Provider, base.SourceControl.Provider).(string),
			Repository: utils.ThisOrThat(opts.SourceControl.Repository, base.SourceControl.Repository).(string),
			Revision:   utils.ThisOrThat(opts.SourceControl.Revision, base.SourceControl.Revision).(string),
		},
		BuilderName:       utils.ThisOrThat(opts.BuilderName, base.BuilderName).(string),
		ReleaseStage:      utils.ThisOrThat(opts.ReleaseStage, base.ReleaseStage).(string),
		AppVersion:        utils.ThisOrThat(opts.AppVersion, base.AppVersion).(string),
		AutoAssignRelease: utils.ThisOrThatBool(opts.AutoAssignRelease, base.AutoAssignRelease),
		MetaData:          utils.ThisOrThat(opts.MetaData, base.MetaData).(map[string]string),
	}
}

func (opts CreateBuildInfo) Validate() error {
	if opts.ApiKey == "" {
		return fmt.Errorf("Missing API Key")
	}

	if opts.AppVersion == "" {
		return fmt.Errorf("Missing App Version")
	}

	if opts.SourceControl.Repository == "" || opts.SourceControl.Revision == "" {
		return fmt.Errorf("Missing Source Control Repository or Revision")
	}

	return nil
}

func PopulateFromCliOpts(opts options.CLI) CreateBuildInfo {
	return CreateBuildInfo{
		ApiKey:           opts.ApiKey,
		AppVersionCode:   opts.CreateBuild.VersionCode,
		AppBundleVersion: opts.CreateBuild.BundleVersion,
		SourceControl: SourceControl{
			Provider:   opts.CreateBuild.Provider,
			Repository: opts.CreateBuild.Repository,
			Revision:   opts.CreateBuild.Revision,
		},
		BuilderName:       opts.CreateBuild.BuilderName,
		ReleaseStage:      opts.CreateBuild.ReleaseStage,
		AppVersion:        opts.CreateBuild.VersionName,
		AutoAssignRelease: &opts.CreateBuild.AutoAssignRelease,
		MetaData:          opts.CreateBuild.Metadata,
	}
}

func PopulateFromPath(path string) CreateBuildInfo {
	return CreateBuildInfo{
		ApiKey:           "",
		AppVersionCode:   "",
		AppBundleVersion: "",
		SourceControl: SourceControl{
			Provider:   "",
			Repository: utils.GetRepoUrl(path),
			Revision:   utils.GetCommitHash(),
		},
		BuilderName:       utils.GetSystemUser(),
		ReleaseStage:      "",
		AppVersion:        "",
		AutoAssignRelease: nil,
		MetaData:          nil,
	}
}

func PopulateFromAndroidManifest(path string) CreateBuildInfo {
	var apiKey string
	androidData, err := android.BuildAndroidInfo(path)

	if err != nil {
		return CreateBuildInfo{}
	}

	for key, value := range androidData.Application.MetaData.Name {
		if value == "com.bugsnag.android.API_KEY" {
			apiKey = androidData.Application.MetaData.Value[key]
		}
	}

	return CreateBuildInfo{
		ApiKey:         apiKey,
		AppVersionCode: androidData.VersionCode,
		AppVersion:     androidData.VersionName,
	}
}

func TheGreatMerge(p1 CreateBuildInfo, p2 CreateBuildInfo) CreateBuildInfo {
	return p1.Override(p2)
}

func BuildInfo(opts options.CLI) (CreateBuildInfo, error) {
	var androidManifestPath string
	var err error
	var BaseOptions CreateBuildInfo

	UserBuildOptions := PopulateFromCliOpts(opts)

	BaseOptions = PopulateFromPath(opts.CreateBuild.Path[0])

	if opts.CreateBuild.AndroidAab != "" {
		androidManifestPath, err = android.GetAndroidManifestFileFromAAB(string(opts.CreateBuild.AndroidAab))

		if err != nil {
			return CreateBuildInfo{}, err
		}
	}

	if androidManifestPath == "" {
		androidManifestPath = string(opts.CreateBuild.AppManifest)
	}

	if androidManifestPath != "" {
		ManifestBuildOptions := PopulateFromAndroidManifest(androidManifestPath)
		BaseOptions = TheGreatMerge(BaseOptions, ManifestBuildOptions)
	}

	FinalMerge := TheGreatMerge(UserBuildOptions, BaseOptions)

	return FinalMerge, nil
}
