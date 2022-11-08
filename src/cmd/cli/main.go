package main

import (
	"flag"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/orchestration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"log"
	"os"
	"strings"
)

var (
	hostArgument    = flag.String("host", "", "Host to search")
	objectArgument  = flag.String("object", "", "Type of object to search")
	includeArgument = flag.String("include", "", "Optional items to include")
	sinceArgument   = flag.String("since", "", "Only search for objects updated since this date")
	orgsArgument    = flag.String("orgs", "", "Comma delimited list of organizations to extract data on")
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
		since := core.ParseDate(sinceArgument)

		includeDetails, includeOwners, includePullRequests := parseIncludeArguments()
		orgsToInclude := parseOrganizations()
		err := orchestration.ExtractRepositories(*hostArgument, since, includeDetails, includeOwners, includePullRequests, orgsToInclude, configurationService, directoryRepository, messageHub)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalln("Unrecognized object")
	}
}

func parseIncludeArguments() (bool, bool, bool) {
	includeDetails := true
	includeOwners := false
	includePullRequests := false

	if *includeArgument == "" {
		includeDetails = true
		includeOwners = false
		includePullRequests = false
	}

	if strings.Contains(*includeArgument, "owner") {
		includeDetails = false
		includeOwners = true
	}

	if strings.Contains(*includeArgument, "pullrequest") {
		includeDetails = false
		includePullRequests = true
	}

	if strings.Contains(*includeArgument, "detail") {
		includeDetails = true
	}
	return includeDetails, includeOwners, includePullRequests
}

func parseOrganizations() []string {
	if strings.TrimSpace(*orgsArgument) == "" {
		return make([]string, 0)
	}
	return strings.Split(*orgsArgument, ",")
}
