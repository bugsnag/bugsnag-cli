package data_access

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

const baseURL = "http://localhost:3000"

type BugsnagOrganization struct {
	ID        string
	AuthToken string
}

type BugsnagProject struct {
	Name                   string         `json:"name"`
	Type                   string         `json:"type"`
	GlobalGrouping         []string       `json:"global_grouping,omitempty"`
	LocationGrouping       []string       `json:"location_grouping,omitempty"`
	DiscardedAppVersions   []string       `json:"discarded_app_versions,omitempty"`
	DiscardedErrors        []string       `json:"discarded_errors,omitempty"`
	URLWhitelist           []string       `json:"url_whitelist,omitempty"`
	IgnoreOldBrowsers      bool           `json:"ignore_old_browsers,omitempty"`
	IgnoredBrowserVersions map[string]any `json:"ignored_browser_versions,omitempty"`
	ResolveOnDeploy        bool           `json:"resolve_on_deploy,omitempty"`
	CollaboratorIDs        []string       `json:"collaborator_ids,omitempty"`
}

// ProjectClient implements project-related operations using a shared Service.
type ProjectClient struct {
	service *Service
}

func (p *ProjectClient) Create(globalOptions options.CLI, logger *log.LoggerWrapper) error {
	if err := p.service.checkEnvironmentVariables(); err != nil {
		return err
	}

	if globalOptions.Create.Project.Name == "" || globalOptions.Create.Project.Type == "" {
		return fmt.Errorf("project name and type must be provided to create a new project")
	}

	project := globalOptions.Create.Project
	createProjectURL := fmt.Sprintf("%s/organizations/%s/projects", baseURL, p.service.bugsnagOrg.ID)
	bugsnagProject := BugsnagProject{Name: project.Name, Type: project.Type, IgnoreOldBrowsers: project.IgnoreOldBrowsers}

	jsonBody, err := json.Marshal(bugsnagProject)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, createProjectURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	responseBody, err := p.service.sendRequest(req)
	if err != nil {
		return err
	}
	logger.Info("Successfully created a new BugSnag project with the following details:\n" + responseBody)

	return nil
}

func (p *ProjectClient) Get(globalOptions options.CLI, logger *log.LoggerWrapper) error {
	if err := p.service.checkEnvironmentVariables(); err != nil {
		return err
	}

	if globalOptions.Get.Project.ID == "" {
		return fmt.Errorf("project ID must be provided to get project details")
	}

	projectID := globalOptions.Get.Project.ID
	getProjectURL := fmt.Sprintf("%s/projects/%s", baseURL, projectID)

	req, err := http.NewRequest(http.MethodGet, getProjectURL, nil)
	if err != nil {
		return err
	}

	responseBody, err := p.service.sendRequest(req)
	if err != nil {
		return err
	}
	logger.Info("Successfully retrieved details for a BugSnag project:\n" + responseBody)

	return nil
}

func (p *ProjectClient) Update(globalOptions options.CLI, logger *log.LoggerWrapper) error {
	if err := p.service.checkEnvironmentVariables(); err != nil {
		return err
	}

	if globalOptions.Update.Project.ID == "" {
		return fmt.Errorf("project ID must be provided to update a project")
	}

	if globalOptions.Update.Project.Name == "" || globalOptions.Update.Project.Type == "" {
		return fmt.Errorf("project name and type must be provided to update a project")
	}

	project := globalOptions.Update.Project

	bugsnagProject := BugsnagProject{
		Name:                   project.Name,
		Type:                   project.Type,
		GlobalGrouping:         project.GlobalGrouping,
		LocationGrouping:       project.LocationGrouping,
		DiscardedAppVersions:   project.DiscardedAppVersions,
		DiscardedErrors:        project.DiscardedErrors,
		URLWhitelist:           project.URLWhitelist,
		IgnoreOldBrowsers:      project.IgnoreOldBrowsers,
		IgnoredBrowserVersions: project.IgnoredBrowserVersions,
		ResolveOnDeploy:        project.ResolveOnDeploy,
		CollaboratorIDs:        project.CollaboratorIDs,
	}

	jsonBody, err := json.Marshal(bugsnagProject)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/projects/%s", baseURL, project.ID), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	responseBody, err := p.service.sendRequest(req)
	if err != nil {
		return err
	}
	logger.Info("Successfully updated BugSnag project with the following details:\n" + responseBody)

	return nil
}

func (s *Service) checkEnvironmentVariables() error {
	if s.bugsnagOrg.ID == "" {
		return fmt.Errorf("missing required environment variable BUGSNAG_ACCOUNT_ID")
	}

	if s.bugsnagOrg.AuthToken == "" {
		return fmt.Errorf("missing required environment variable BUGSNAG_CLI_PERSONAL_AUTH_TOKEN")
	}

	return nil
}

func (s *Service) sendRequest(req *http.Request) (string, error) {
	req.Header.Set("X-Bugsnag-Api", "true")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", s.bugsnagOrg.AuthToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
