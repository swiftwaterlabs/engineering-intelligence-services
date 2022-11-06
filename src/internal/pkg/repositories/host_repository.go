package repositories

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type HostRepository interface {
	GetAll() ([]*models.Host, error)
	Get(identifier string) (*models.Host, error)
}

func NewHostRepository(appConfig *configuration.AppConfig, config configuration.ConfigurationService) HostRepository {
	return NewDynamoDbHostRepository(appConfig, config)
}
