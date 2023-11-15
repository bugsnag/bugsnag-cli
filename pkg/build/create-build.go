package build

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

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
	AutoAssignRelease bool
	MetaData          map[string]string
}

func (opts CreateBuildInfo) Override(base CreateBuildInfo) CreateBuildInfo {
	return CreateBuildInfo{
		GeneralInfo: GeneralInfo{
			ApiKey: utils.Xor(opts.GeneralInfo.ApiKey, base.GeneralInfo.ApiKey).(string),
		},
		AndroidInfo: AndroidInfo{
			AppVersionCode: utils.Xor(opts.AndroidInfo.AppVersionCode, base.AndroidInfo.AppVersionCode).(string),
		},
		IosInfo: IosInfo{
			AppBundleVersion: utils.Xor(opts.IosInfo.AppBundleVersion, base.IosInfo.AppBundleVersion).(string),
		},
		SourceControl: SourceControl{
			Provider:   utils.Xor(opts.SourceControl.Provider, base.SourceControl.Provider).(string),
			Repository: utils.Xor(opts.SourceControl.Repository, base.SourceControl.Repository).(string),
			Revision:   utils.Xor(opts.SourceControl.Revision, base.SourceControl.Revision).(string),
		},
		BuilderName:       utils.Xor(opts.BuilderName, base.BuilderName).(string),
		ReleaseStage:      utils.Xor(opts.ReleaseStage, base.ReleaseStage).(string),
		AppVersion:        utils.Xor(opts.AppVersion, base.AppVersion).(string),
		AutoAssignRelease: utils.Xor(opts.AutoAssignRelease, base.AutoAssignRelease).(bool),
		MetaData:          utils.Xor(opts.MetaData, base.MetaData).(map[string]string),
	}
}
