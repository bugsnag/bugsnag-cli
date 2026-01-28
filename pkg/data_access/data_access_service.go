package data_access

import (
	"net/http"
	"os"
	"time"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// Projects defines the project-related behaviour exposed by the data access layer.
type Projects interface {
	Create(globalOptions options.CLI, logger *log.LoggerWrapper) error
	Get(globalOptions options.CLI, logger *log.LoggerWrapper) error
	Update(globalOptions options.CLI, logger *log.LoggerWrapper) error
}

type Service struct {
	httpClient *http.Client
	bugsnagOrg BugsnagOrganization

	// Projects provides access to project-related operations, e.g. dataAccessService.Projects.Create(...).
	Projects Projects
}

func NewService() *Service {
	s := &Service{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		bugsnagOrg: BugsnagOrganization{
			ID:        os.Getenv("BUGSNAG_ACCOUNT_ID"),
			AuthToken: os.Getenv("BUGSNAG_CLI_PERSONAL_AUTH_TOKEN"),
		},
	}

	// Wire a concrete ProjectClient (defined in projects.go) that delegates to this Service's shared state.
	s.Projects = &ProjectClient{service: s}

	return s
}
