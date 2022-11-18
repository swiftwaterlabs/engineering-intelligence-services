package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/orchestration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"log"
	"os"
	"strings"
)

var (
	appConfig             *configuration.AppConfig
	configurationService  configuration.ConfigurationService
	messageHub            messaging.MessageHub
	eventSourceRepository repositories.EventSourceRepository
)

func init() {
	appConfig = &configuration.AppConfig{
		AwsRegion:          os.Getenv("aws_region"),
		AuthenticateEvents: os.Getenv("authenticate_events"),
	}
	configurationService = configuration.NewConfigurationService(appConfig)
	messageHub = messaging.NewMessageHub(appConfig)
	eventSourceRepository = repositories.NewEventSourceRepository(appConfig, configurationService)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authenticateEvents := strings.EqualFold(appConfig.AuthenticateEvents, "true")

	isAuthenticated, err := orchestration.ProcessWebhookEvent(event.Headers, event.Body, authenticateEvents, eventSourceRepository, configurationService, messageHub)
	if err != nil {
		log.Printf("error when processing webhoook|%s", err)
	}

	if !isAuthenticated {
		return events.APIGatewayProxyResponse{StatusCode: 403}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
