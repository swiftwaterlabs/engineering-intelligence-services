package configuration

type ConfigurationService interface {
	GetValue(name string) string
	GetSecret(name string) string
}

func NewConfigurationService(appConfig *AppConfig) ConfigurationService {
	return NewSecretsManagerConfigurationService(appConfig)
}
