package clients

import (
	"fmt"
	"github.com/swiftwaterlabs/engineering-intelligence-services/external/sonargo"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
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

	// Note: processProjects is preferred, but requires admin rights.
	err = c.processComponents(client, options, handler)

	return err
}

func (c *SonarqubeTestResultClient) getClient(configurationService configuration.ConfigurationService) (*sonargo.Client, error) {
	clientSecret := configurationService.GetSecret(c.host.ClientSecretName)
	client, err := sonargo.NewClientByToken(c.host.BaseUrl, clientSecret)

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
		if projects == nil {
			break
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

func (c *SonarqubeTestResultClient) processComponents(client *sonargo.Client,
	options *models.TestResultProcessingOptions,
	handler func(data []*models.TestResult)) error {

	searchOptions := &sonargo.ComponentsSearchOption{
		Q:          strings.Join(options.Projects, ","),
		Ps:         "500",
		Qualifiers: sonargo.QualifierProject,
	}

	processingErrors := make([]error, 0)

	currentPage := 1
	const maxNumberOfPages = 20
	for {
		components, _, err := client.Components.Search(searchOptions)
		if err != nil {
			processingErrors = append(processingErrors, err)
		}
		if components == nil {
			break
		}
		if !c.hasComponentData(components) {
			break
		}

		for _, item := range components.Components {
			err = c.processComponent(client, item, options, handler)
			if err != nil {
				processingErrors = append(processingErrors, err)
			}
		}

		currentPage++
		if currentPage > maxNumberOfPages {
			break
		}
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

	log.Printf("Processing test results for %s", component.Project)
	searchOptions := &sonargo.MeasuresSearchHistoryOption{
		Component: component.Key,
		From:      "",
		Metrics:   "coverage,new_coverage,new_uncovered_lines,lines_to_cover,uncovered_lines",
		P:         "",
		Ps:        "100",
	}

	if options.Since != nil {
		searchOptions.From = options.Since.Format("2006-01-02")
	}

	processingErrors := make([]error, 0)

	currentPage := 1
	for {
		measuresData, _, err := client.Measures.SearchHistory(searchOptions)
		if err != nil {
			processingErrors = append(processingErrors, err)
		}

		if !c.hasMeasuresData(measuresData) {
			break
		}

		testResults := mapTestResult(c.host, component, measuresData)
		handler(testResults)

		currentPage++
		searchOptions.P = fmt.Sprint(currentPage)
	}

	if len(processingErrors) == 0 {
		return nil
	}
	return core.ConsolidateErrors(processingErrors)
}

func (c *SonarqubeTestResultClient) hasComponentData(data *sonargo.ComponentsSearchObject) bool {
	if data == nil {
		return false
	}

	if len(data.Components) > 1 {
		return true
	}

	return false
}

func (c *SonarqubeTestResultClient) hasMeasuresData(data *sonargo.MeasuresSearchHistoryObject) bool {
	if data == nil {
		return false
	}

	for _, measure := range data.Measures {
		if len(measure.Histories) > 1 {
			return true
		}
	}

	return false
}
