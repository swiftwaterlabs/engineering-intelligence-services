package clients

import (
	"errors"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"strings"
)

type SourceCodeRepositoryClient interface {
	ProcessRepositories(configurationService configuration.ConfigurationService,
		options *models.RepositoryProcessingOptions,
		handlers *RepositoryHandlers) error
}

func NewSourceCodeRepositoryClient(host *models.Host) (SourceCodeRepositoryClient, error) {
	if strings.Contains(strings.ToLower(host.SubType), "github") {
		return &GithubSourceCodeRepositoryClient{
			host: host,
		}, nil
	}

	return nil, errors.New("unrecognized host type")
}

type RepositoryHandlers struct {
	Repository  func(data []*models.Repository)
	Owner       func(data []*models.RepositoryOwner)
	PullRequest func(data []*models.PullRequest)
	BranchRule  func(data []*models.BranchProtectionRule)
}
