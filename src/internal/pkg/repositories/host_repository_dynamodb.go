package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"strings"
)

func NewDynamoDbHostRepository(appConfig *configuration.AppConfig, config configuration.ConfigurationService) *DynamoDbDirectoryRepository {
	session := configuration.GetAwsSession(appConfig)
	client := dynamodb.New(session)

	return &DynamoDbDirectoryRepository{
		tableName: config.GetValue("engineering_intelligence_prd_directories_table"),
		client:    client,
	}
}

type DynamoDbDirectoryRepository struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func (r *DynamoDbDirectoryRepository) GetAll(hostType string) ([]*models.Host, error) {
	result := make([]*models.Host, 0)
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}
	queryResult, err := r.client.Scan(scanInput)
	if err != nil {
		return result, err
	}

	for _, item := range queryResult.Items {
		host := r.mapItemToHost(item)
		if strings.EqualFold(hostType, host.Type) {
			result = append(result, host)
		}
	}

	return result, nil
}

func (r *DynamoDbDirectoryRepository) Get(identifier string) (*models.Host, error) {
	itemInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(identifier),
			},
		},
		TableName: aws.String(r.tableName),
	}
	queryResult, err := r.client.GetItem(itemInput)
	if err != nil {
		return nil, err
	}

	result := r.mapItemToHost(queryResult.Item)
	return result, nil
}

func (r *DynamoDbDirectoryRepository) mapItemToHost(item map[string]*dynamodb.AttributeValue) *models.Host {
	return &models.Host{
		Id:                 r.getStringValue(item["Id"]),
		Name:               r.getStringValue(item["Name"]),
		BaseUrl:            r.getStringValue(item["BaseUrl"]),
		Type:               r.getStringValue(item["Type"]),
		SubType:            r.getStringValue(item["SubType"]),
		AuthenticationType: r.getStringValue(item["AuthenticationType"]),
		ClientSecretName:   r.getStringValue(item["ClientSecretName"]),
	}
}

func (r *DynamoDbDirectoryRepository) getStringValue(item *dynamodb.AttributeValue) string {
	if item == nil || item.S == nil {
		return ""
	}
	return aws.StringValue(item.S)
}
