package clients

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type SonarqubeTestResultClient struct {
	host *models.Host
}

func (c *SonarqubeTestResultClient) ProcessTestResults(configurationService configuration.ConfigurationService,
	options *models.TestResultProcessingOptions,
	handlers func(data []*models.TestResult)) error {
	return nil
}
