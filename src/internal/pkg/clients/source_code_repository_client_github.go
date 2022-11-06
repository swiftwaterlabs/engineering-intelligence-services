package clients

import (
	"context"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
)

type GithubSourceCodeRepositoryClient struct {
	host *models.Host
}

func (c *GithubSourceCodeRepositoryClient) ProcessRepositories(configurationService configuration.ConfigurationService,
	processor func(data []*models.Repository)) error {

	hostSecret := configurationService.GetSecret(c.host.ClientSecretName)
	client, err := GetGitHubClient(c.host.SubType, c.host.BaseUrl, c.host.AuthenticationType, hostSecret)
	if err != nil {
		return err
	}

	user, response, err := client.Users.Get(context.Background(), "jrolstad")
	if err != nil {
		return err
	}
	log.Printf("Users:%s|Response Code:%v", user.GetURL(), response.StatusCode)

	return nil
}
