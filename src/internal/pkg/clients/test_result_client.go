package clients

import (
	"errors"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"strings"
)

type TestResultClient interface {
	ProcessTestResults(configurationService configuration.ConfigurationService,
		options *models.TestResultProcessingOptions,
		handlers func(data []*models.TestResult)) error
}

func NewTestResultClient(host *models.Host) (TestResultClient, error) {
	if strings.Contains(strings.ToLower(host.SubType), "sonarqube") {
		return &SonarqubeTestResultClient{
			host: host,
		}, nil
	}

	return nil, errors.New("unrecognized host type")
}
