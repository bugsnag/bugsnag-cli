package data_access

type Project struct {
	ID                     string         `help:"The ID of the project to retrieve"`
	Name                   string         `help:"The name of the project"`
	Type                   string         `help:"The project type (e.g. 'js', 'android', 'ios')"`
	GlobalGrouping         []string       `help:"List of error classes to group globally by class (global_grouping)"`
	LocationGrouping       []string       `help:"List of error classes to group by context (location_grouping)"`
	DiscardedAppVersions   []string       `help:"App versions to discard events for (discarded_app_versions). Supports regex and semver ranges"`
	DiscardedErrors        []string       `help:"Error classes to discard events for (discarded_errors)"`
	URLWhitelist           []string       `help:"List of whitelisted script source domains (url_whitelist)"`
	IgnoreOldBrowsers      bool           `help:"Whether to ignore old browsers for this project"`
	IgnoredBrowserVersions map[string]any `help:"Ignored browser versions for JS projects (ignored_browser_versions)"`
	ResolveOnDeploy        bool           `help:"Whether to mark all errors as fixed on deploy (resolve_on_deploy)"`
	CollaboratorIDs        []string       `help:"List of collaborator IDs to set on the project (collaborator_ids)"`
}

type ProjectActions struct {
	Create createProject `cmd:"" help:"Create a project"`
	Get    getProject    `cmd:"" help:"Get a project"`
}

type createProject struct {
	Name              string `help:"The name of the project" required:""`
	Type              string `help:"The project type (e.g. 'js', 'android', 'ios')" required:""`
	IgnoreOldBrowsers bool   `help:"Whether to ignore old browsers for this project" default:"true"`
}

type getProject struct {
	ProjectID string `help:"The ID of the project to retrieve" required:""`
}
