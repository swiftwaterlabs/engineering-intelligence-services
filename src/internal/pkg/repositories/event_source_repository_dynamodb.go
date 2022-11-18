package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type DynamoDbEventSourceRepository struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func (r *DynamoDbEventSourceRepository) GetAllActive() ([]*models.EventSource, error) {
	result := make([]*models.EventSource, 0)
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}
	queryResult, err := r.client.Scan(scanInput)
	if err != nil {
		return result, err
	}

	for _, item := range queryResult.Items {
		eventSource := r.mapItemEventSource(item)
		if eventSource.Active {
			result = append(result, eventSource)
		}
	}

	return result, nil
}

func (r *DynamoDbEventSourceRepository) mapItemEventSource(item map[string]*dynamodb.AttributeValue) *models.EventSource {
	return &models.EventSource{
		Id:      getStringValue(item["Id"]),
		Name:    getStringValue(item["Name"]),
		Type:    getStringValue(item["Type"]),
		Active:  getBooleanValue(item["Active"]),
		Secrets: getAuthenticationSecretValues(item["Secrets"]),
	}
}

func getAuthenticationSecretValues(item *dynamodb.AttributeValue) []*models.AuthenticationSecret {
	if item == nil || item.L == nil {
		return make([]*models.AuthenticationSecret, 0)
	}

	results := make([]*models.AuthenticationSecret, 0)
	for _, value := range item.L {
		if value.M != nil {
			secret := &models.AuthenticationSecret{
				SecretName: getStringValue(value.M["SecretName"]),
				Active:     getBooleanValue(value.M["Active"]),
			}
			results = append(results, secret)
		}
	}

	return results
}
