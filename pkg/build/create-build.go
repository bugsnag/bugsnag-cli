package build

import (
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/carlmjohnson/truthy"
)

type GeneralInfo struct {
	ApiKey string
}

type AndroidInfo struct {
	AppVersionCode string
}

type IosInfo struct {
	AppBundleVersion string
}

type SourceControl struct {
	Provider   string
	Repository string
	Revision   string
}

type CreateBuildInfo struct {
	GeneralInfo       GeneralInfo
	AndroidInfo       AndroidInfo
	IosInfo           IosInfo
	SourceControl     SourceControl
	BuilderName       string
	ReleaseStage      string
	AppVersion        string
	AutoAssignRelease *bool
	MetaData          map[string]string
}

func (opts CreateBuildInfo) Override(base CreateBuildInfo) CreateBuildInfo {
	var apiKey string

	if truthy.Value(opts.GeneralInfo.ApiKey) {
		apiKey = opts.GeneralInfo.ApiKey
	} else {
		apiKey = base.GeneralInfo.ApiKey
	}

	return CreateBuildInfo{
		GeneralInfo: GeneralInfo{
			ApiKey: apiKey,
		},
		AndroidInfo: AndroidInfo{
			AppVersionCode: utils.ThisOrThat(opts.AndroidInfo.AppVersionCode, base.AndroidInfo.AppVersionCode).(string),
		},
		IosInfo: IosInfo{
			AppBundleVersion: utils.ThisOrThat(opts.IosInfo.AppBundleVersion, base.IosInfo.AppBundleVersion).(string),
		},
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

func PopulateFromCliOpts(opts options.CLI) CreateBuildInfo {
	return CreateBuildInfo{
		GeneralInfo: GeneralInfo{ApiKey: opts.ApiKey},
		AndroidInfo: AndroidInfo{AppVersionCode: opts.CreateBuild.VersionCode},
		IosInfo:     IosInfo{AppBundleVersion: opts.CreateBuild.BundleVersion},
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
		GeneralInfo: GeneralInfo{ApiKey: "foobar"},
		AndroidInfo: AndroidInfo{AppVersionCode: ""},
		IosInfo:     IosInfo{AppBundleVersion: ""},
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

func TheGreatMerge(p1 CreateBuildInfo, p2 CreateBuildInfo) CreateBuildInfo {
	return p1.Override(p2)
}
