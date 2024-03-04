package ios

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

func FindProjectOrWorkspaceInPath(path string) (string, error) {
	var xcodeProjectPath string
	var err error
	xcodeProjectPath, err = utils.FindLatestFileWithSuffix(path, ".xcodeproj")

	if err != nil {
		xcodeProjectPath, err = utils.FindLatestFileWithSuffix(path, ".xcworkspace")
		if err != nil {
			return "", err
		}
	}

	return xcodeProjectPath, nil
}
