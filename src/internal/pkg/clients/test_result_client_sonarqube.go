package clients

import (
	"fmt"
	sonargo "github.com/magicsong/sonargo/sonar"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"strings"
)

type SonarqubeTestResultClient struct {
	host *models.Host
}

func (c *SonarqubeTestResultClient) ProcessTestResults(configurationService configuration.ConfigurationService,
	options *models.TestResultProcessingOptions,
	handler func(data []*models.TestResult)) error {

	client, err := c.getClient(configurationService)
	if err != nil {
		return err
	}

	err = c.processProjects(client, options, handler)

	return err
}

func (c *SonarqubeTestResultClient) getClient(configurationService configuration.ConfigurationService) (*sonargo.Client, error) {
	clientSecret := configurationService.GetSecret(c.host.ClientSecretName)
	client, err := sonargo.NewClient(c.host.BaseUrl, c.host.ClientId, clientSecret)

	return client, err
}

func (c *SonarqubeTestResultClient) processProjects(client *sonargo.Client,
	options *models.TestResultProcessingOptions,
	handler func(data []*models.TestResult)) error {

	searchOptions := &sonargo.ProjectsSearchOption{
		Projects:   strings.Join(options.Projects, ","),
		Ps:         "100",
		Qualifiers: "",
	}

	processingErrors := make([]error, 0)

	currentPage := 1
	for {
		projects, _, err := client.Projects.Search(searchOptions)
		if err != nil {
			processingErrors = append(processingErrors, err)
		}

		for _, item := range projects.Components {
			err = c.processComponent(client, item, options, handler)
			if err != nil {
				processingErrors = append(processingErrors, err)
			}
		}

		if currentPage >= projects.Paging.Total {
			break
		}
		currentPage++
		searchOptions.P = fmt.Sprint(currentPage)
	}

	if len(processingErrors) == 0 {
		return nil
	}
	return core.ConsolidateErrors(processingErrors)
}

func (c *SonarqubeTestResultClient) processComponent(client *sonargo.Client,
	component *sonargo.Component,
	options *models.TestResultProcessingOptions,
	handler func(data []*models.TestResult)) error {

	searchOptions := &sonargo.MeasuresSearchHistoryOption{
		Component: component.Key,
		From:      "",
		Metrics:   "overage,new_coverage,new_uncovered_lines",
		P:         "",
		Ps:        "100",
	}

	if options.Since != nil {
		searchOptions.From = options.Since.Format("2006-02-01T15:04:05")
	}

	processingErrors := make([]error, 0)

	currentPage := 1
	for {
		measuresData, _, err := client.Measures.SearchHistory(searchOptions)
		if err != nil {
			processingErrors = append(processingErrors, err)
		}

		measuresByDate := make(map[string]map[string]string, 0)
		for _, measure := range measuresData.Measures {
			for _, history := range measure.Histories {
				if measuresByDate[history.Date] == nil {
					measuresByDate[history.Date] = make(map[string]string, 0)
				}
				measuresByDate[history.Date][measure.Metric] = history.Value
			}
		}

		testResults := mapTestResult(component, measuresData)
		handler(testResults)

		if currentPage >= measuresData.Paging.Total {
			break
		}
		currentPage++
		searchOptions.P = fmt.Sprint(currentPage)
	}

	if len(processingErrors) == 0 {
		return nil
	}
	return core.ConsolidateErrors(processingErrors)
}
