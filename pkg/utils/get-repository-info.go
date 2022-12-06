package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetRepoUrl - Gets the URl of a git repo.
func GetRepoUrl() (string, error) {
	gitLocation, err := exec.LookPath("git")

	if err != nil {
		return "", fmt.Errorf("unable to find git on system: %w", err)
	}

	cmd := exec.Command(gitLocation, "config", "--get", "remote.origin.url")

	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	url := ParseGitUrl(strings.TrimSuffix(string(cmdOutput), "\n"))

	return url, nil
}

//ParseGitUrl - Parses a git url to ensure that it is in HTTPS format
func ParseGitUrl(url string) string {
	if strings.Contains(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
		url = strings.Replace(url, ".git", "", 1)
		return url
	}
	return url
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
