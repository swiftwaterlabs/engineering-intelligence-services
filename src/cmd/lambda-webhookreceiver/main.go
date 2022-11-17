package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/orchestration"
	"log"
	"os"
)

var (
	configurationService configuration.ConfigurationService
	messageHub           messaging.MessageHub
)

func init() {
	appConfig := &configuration.AppConfig{
		AwsRegion: os.Getenv("aws_region"),
	}
	configurationService = configuration.NewConfigurationService(appConfig)
	messageHub = messaging.NewMessageHub(appConfig)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	err := orchestration.ProcessWebhookEvent(event.Headers, event.Body, configurationService, messageHub)
	if err != nil {
		log.Printf("error when processing webhoook|%s", err)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil

}
