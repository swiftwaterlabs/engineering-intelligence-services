package clients

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
)

func (c *GithubSourceCodeRepositoryClient) processWebhookForOrganization(client *github.Client,
	organization *models.Organization) []*models.Webhook {

	options := &github.ListOptions{
		Page:    0,
		PerPage: 100,
	}
	result := make([]*models.Webhook, 0)

	for {
		data, response, err := client.Organizations.ListHooks(context.Background(), organization.Name, options)
		if err != nil {
			log.Printf("Unable to retrieve webhooks for %s", organization.Url)
		}

		for _, item := range data {
			mappedItem := mapWebHook(organization, nil, item)
			result = append(result, mappedItem)
		}

		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	return result
}

func (c *GithubSourceCodeRepositoryClient) processWebhookForRepository(client *github.Client,
	repository *models.Repository) []*models.Webhook {
	options := &github.ListOptions{
		Page:    0,
		PerPage: 100,
	}
	result := make([]*models.Webhook, 0)

	for {
		data, response, err := client.Repositories.ListHooks(context.Background(), repository.Organization.Name, repository.Name, options)
		if err != nil {
			log.Printf("Unable to retrieve webhooks for %s", repository.Url)
		}

		for _, item := range data {
			mappedItem := mapWebHook(nil, repository, item)
			result = append(result, mappedItem)
		}

		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	return result
}
