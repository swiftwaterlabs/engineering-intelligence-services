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
	hostArgument   = flag.String("host", "", "Host to search")
	objectArgument = flag.String("object", "", "Type of object to search")
	sinceArgument  = flag.String("since", "", "Only search for objects updated since this date")
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
		err := orchestration.ExtractRepositories(*hostArgument, since, configurationService, directoryRepository, messageHub)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalln("Unrecognized object")
	}
}
