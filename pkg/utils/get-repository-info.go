package utils

import (
	"os/exec"
	"strings"
)

// GetRepoUrl - Gets the URl of a git repo.
func GetRepoUrl(repoPath string) string {
	gitLocation, err := exec.LookPath("git")

	if err != nil {
		return ""
	}

	remoteOriginCmd := exec.Command(gitLocation, "-C", repoPath, "config", "--get", "remote.origin.url")
	remoteOriginCmdOutput, err := remoteOriginCmd.CombinedOutput()

	if err != nil {
		remoteCmd := exec.Command(gitLocation, "-C", repoPath, "remote")
		remoteCmdOutput, err := remoteCmd.CombinedOutput()
		if err != nil {
			return ""
		}
		remotes := strings.Split(string(remoteCmdOutput), "\n")
		remoteOriginCmd = exec.Command(gitLocation, "-C", repoPath, "config", "--get", "remote."+remotes[0]+".url")
		remoteOriginCmdOutput, err = remoteOriginCmd.CombinedOutput()
		if err != nil {
			return ""
		}
	}

	return strings.TrimSuffix(string(remoteOriginCmdOutput), "\n")
}

// GetCommitHash - Gets the commit hash from a repo
func GetCommitHash() string {
	gitLocation, err := exec.LookPath("git")

	if err != nil {
		return ""
	}

	cmd := exec.Command(gitLocation, "rev-parse", "HEAD")

	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		return ""
	}

	return strings.TrimSuffix(string(cmdOutput), "\n")
}
