package repositories

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type EventSourceRepository interface {
	GetAllActive() ([]*models.EventSource, error)
}

func NewEventSourceRepository(appConfig *configuration.AppConfig, config configuration.ConfigurationService) EventSourceRepository {
	session := configuration.GetAwsSession(appConfig)
	client := dynamodb.New(session)

	return &DynamoDbEventSourceRepository{
		tableName: config.GetValue("engineering_intelligence_prd_event_sources_table"),
		client:    client,
	}
}
