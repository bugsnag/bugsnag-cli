package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetRepoUrl - Gets the URl of a git repo.
func GetRepoUrl() string {
	gitLocation, err := exec.LookPath("git")

	if err != nil {
		return ""
	}

	remoteOriginCmd := exec.Command(gitLocation, "config", "--get", "remote.origin.url")
	remoteOriginCmdOutput, err := remoteOriginCmd.CombinedOutput()

	if err != nil {
		remoteCmd := exec.Command(gitLocation, "remote")
		remoteCmdOutput, err := remoteCmd.CombinedOutput()
		if err != nil {
			return ""
		}
		remotes := strings.Split(string(remoteCmdOutput), "\n")
		remoteOriginCmd = exec.Command(gitLocation, "config", "--get", "remote."+remotes[0]+".url")
		remoteOriginCmdOutput, err = remoteOriginCmd.CombinedOutput()
		if err != nil {
			return ""
		}
	}

	return string(strings.TrimSuffix(string(remoteOriginCmdOutput), "\n"))
}

// GetCommitHash - Gets the commit hash from a repo
func GetCommitHash() (string, error) {
	gitLocation, err := exec.LookPath("git")

	if err != nil {
		return "", fmt.Errorf("unable to find git on system: %w", err)
	}

	cmd := exec.Command(gitLocation, "rev-parse", "HEAD")

	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(cmdOutput), "\n"), nil
}
