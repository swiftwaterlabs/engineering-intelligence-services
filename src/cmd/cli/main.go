package main

import (
	"flag"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/orchestration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"log"
	"os"
	"strings"
)

var (
	hostArgument     = flag.String("host", "", "Host to search")
	objectArgument   = flag.String("object", "", "Type of object to search")
	includeArgument  = flag.String("include", "", "Optional items to include")
	sinceArgument    = flag.String("since", "", "Only search for objects updated since this date")
	orgsArgument     = flag.String("orgs", "", "Comma delimited list of organizations to extract data on")
	projectsArgument = flag.String("projects", "", "Comma delimited list of projects to extract data on")
)

func main() {
	flag.Parse()

	appConfig := &configuration.AppConfig{
		AwsRegion: os.Getenv("aws_region"),
	}
	configurationService := configuration.NewConfigurationService(appConfig)
	directoryRepository := repositories.NewHostRepository(appConfig, configurationService)
	messageHub := messaging.NewMessageHub(appConfig)

	switch strings.ToLower(*objectArgument) {
	case "repository":
		options := parseRepositoryArguments()
		err := orchestration.ExtractRepositories(*hostArgument, options, configurationService, directoryRepository, messageHub)
		if err != nil {
			log.Fatal(err)
		}
	case "testresult":
		options := parseTestResultArguments()
		err := orchestration.ExtractTestResults(*hostArgument, options, configurationService, directoryRepository, messageHub)
		if err != nil {
			log.Fatal(err)
		}
	case "webhook":
		{
			err := orchestration.ListenForWebhookEvents(configurationService, directoryRepository, messageHub)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		log.Fatalln("Unrecognized object")
	}
}

func parseRepositoryArguments() *models.RepositoryProcessingOptions {
	result := &models.RepositoryProcessingOptions{
		IncludeDetails:      true,
		IncludeOwners:       false,
		IncludePullRequests: false,
		IncludeBranchRules:  false,
		IncludeWebhooks:     false,
		Organizations:       make([]string, 0),
		Since:               core.ParseDate(sinceArgument),
	}

	if strings.Contains(*includeArgument, "owner") {
		result.IncludeDetails = false
		result.IncludeOwners = true
	}

	if strings.Contains(*includeArgument, "pullrequest") {
		result.IncludeDetails = false
		result.IncludePullRequests = true
	}

	if strings.Contains(*includeArgument, "branchrule") {
		result.IncludeDetails = false
		result.IncludeBranchRules = true
	}

	if strings.Contains(*includeArgument, "hook") {
		result.IncludeDetails = false
		result.IncludeWebhooks = true
	}

	if strings.Contains(*includeArgument, "detail") {
		result.IncludeDetails = true
	}

	if strings.TrimSpace(*orgsArgument) != "" {
		result.Organizations = strings.Split(*orgsArgument, ",")
	}

	return result
}

func parseTestResultArguments() *models.TestResultProcessingOptions {
	result := &models.TestResultProcessingOptions{
		Since:    core.ParseDate(sinceArgument),
		Projects: nil,
	}

	if strings.TrimSpace(*projectsArgument) != "" {
		result.Projects = strings.Split(*projectsArgument, ",")
	}

	return result
}
