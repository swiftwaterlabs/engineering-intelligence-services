package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

func NewDynamoDbHostRepository(appConfig *configuration.AppConfig, config configuration.ConfigurationService) *DynamoDbDirectoryRepository {
	session := configuration.GetAwsSession(appConfig)
	client := dynamodb.New(session)

	return &DynamoDbDirectoryRepository{
		tableName: config.GetValue("identity_intelligence_prd_directories_table"),
		client:    client,
	}
}

type DynamoDbDirectoryRepository struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func (r *DynamoDbDirectoryRepository) GetAll() ([]*models.Host, error) {
	result := make([]*models.Host, 0)
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}
	queryResult, err := r.client.Scan(scanInput)
	if err != nil {
		return result, err
	}

	for _, item := range queryResult.Items {
		directory := r.mapItemToDirectory(item)
		result = append(result, directory)
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

	result := r.mapItemToDirectory(queryResult.Item)
	return result, nil
}

func (r *DynamoDbDirectoryRepository) mapItemToDirectory(item map[string]*dynamodb.AttributeValue) *models.Host {
	return &models.Host{
		Id:                 r.getStringValue(item["Id"]),
		Name:               r.getStringValue(item["Name"]),
		BaseUrl:            r.getStringValue(item["BaseUrl"]),
		Type:               r.getStringValue(item["Type"]),
		AuthenticationType: r.getStringValue(item["AuthenticationType"]),
		ClientSecretName:   r.getStringValue(item["ClientSecretName"]),
	}
}

func (r *DynamoDbDirectoryRepository) getStringValue(item *dynamodb.AttributeValue) string {
	if item.S == nil {
		return ""
	}
	return aws.StringValue(item.S)
}
