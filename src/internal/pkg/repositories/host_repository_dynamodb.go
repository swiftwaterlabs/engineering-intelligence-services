package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
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
		Id:                 getStringValue(item["Id"]),
		Name:               getStringValue(item["Name"]),
		BaseUrl:            getStringValue(item["BaseUrl"]),
		Type:               getStringValue(item["Type"]),
		SubType:            getStringValue(item["SubType"]),
		AuthenticationType: getStringValue(item["AuthenticationType"]),
		ClientSecretName:   getStringValue(item["ClientSecretName"]),
		Options:            getMapValue(item["Options"]),
	}
}

func getStringValue(item *dynamodb.AttributeValue) string {
	if item == nil || item.S == nil {
		return ""
	}
	return aws.StringValue(item.S)
}

func getMapValue(item *dynamodb.AttributeValue) map[string]string {
	if item == nil || item.M == nil {
		return make(map[string]string, 0)
	}

	result := make(map[string]string, 0)

	for key, value := range item.M {
		result[key] = aws.StringValue(value.S)
	}
	log.Println(result)
	return result
}

func getBooleanValue(item *dynamodb.AttributeValue) bool {
	if item == nil || item.BOOL == nil {
		return false
	}
	return aws.BoolValue(item.BOOL)
}
