package configuration

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func NewSecretsManagerConfigurationService(appConfig *AppConfig) ConfigurationService {
	session := GetAwsSession(appConfig)
	client := secretsmanager.New(session)

	return &secretsManagerConfigurationService{
		client: client,
	}
}

type secretsManagerConfigurationService struct {
	client *secretsmanager.SecretsManager
}

func (s *secretsManagerConfigurationService) GetValue(name string) string {
	return s.GetSecret(name)
}

func (s *secretsManagerConfigurationService) GetSecret(name string) string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	secretValue, err := s.client.GetSecretValue(input)
	if err != nil {
		return ""
	}

	return *secretValue.SecretString
}
