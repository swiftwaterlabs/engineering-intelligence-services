package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"os"
)

var (
	repository repositories.EngineeringObjectRepository
)

func init() {
	appConfig := &configuration.AppConfig{
		AwsRegion: os.Getenv("aws_region"),
	}
	configurationService := configuration.NewConfigurationService(appConfig)
	repository = repositories.NewEngineeringObjectRepository(configurationService)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) error {

	for _, item := range event.Records {
		err := processEvent(item)

		if err != nil {
			return err
		}
	}

	return nil
}

func processEvent(message events.SQSMessage) error {
	directoryObject := &models.EngineeringObject{}
	core.MapFromJson(message.Body, directoryObject)
	directoryObject.Data = core.SanitizeString(message.Body)

	err := repository.Save(directoryObject)
	return err
}
