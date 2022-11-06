package clients

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type SourceCodeRepositoryClient interface {
	ProcessRepositories(host *models.Host,
		configurationService configuration.ConfigurationService,
		processor func(data []*models.Repository))
}

func NewSourceCodeRepositoryClient(host *models.Host){

}