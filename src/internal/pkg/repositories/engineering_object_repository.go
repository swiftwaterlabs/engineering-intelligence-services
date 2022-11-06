package repositories

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type EngineeringObjectRepository interface {
	Save(item *models.EngineeringObject) error
	Destroy()
}

func NewEngineeringObjectRepository(config configuration.ConfigurationService) EngineeringObjectRepository {
	instance := &S3EngineeringObjectRepository{}
	instance.init(config)

	return instance
}
