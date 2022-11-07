package clients

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"golang.org/x/oauth2"
	"strings"
)

type GithubSourceCodeRepositoryClient struct {
	host *models.Host
}

const (
	githubClientTypeEnterpriseServer = "GitHub Enterprise Server"
)

func (c *GithubSourceCodeRepositoryClient) ProcessRepositories(configurationService configuration.ConfigurationService,
	repositoryHandler func(data []*models.Repository)) error {

	hostSecret := configurationService.GetSecret(c.host.ClientSecretName)
	client, err := getGitHubClient(c.host.SubType, c.host.BaseUrl, c.host.AuthenticationType, hostSecret)
	if err != nil {
		return err
	}

	if strings.EqualFold(githubClientTypeEnterpriseServer, c.host.SubType) {
		return c.processAllOrganizationsOnHost(client, repositoryHandler)
	}

	return c.processAllMemberOrganizations(client, repositoryHandler)
}

func getGitHubClient(hostType string, baseUrl, authenticationType string, authenticationSecret string) (*github.Client, error) {
	context := context.Background()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authenticationSecret},
	)
	tokenClient := oauth2.NewClient(context, tokenSource)

	if strings.EqualFold(githubClientTypeEnterpriseServer, hostType) {
		client, err := github.NewEnterpriseClient(baseUrl, baseUrl, tokenClient)
		return client, err
	}

	client := github.NewClient(tokenClient)
	return client, nil
}

func (c *GithubSourceCodeRepositoryClient) processAllOrganizationsOnHost(client *github.Client,
	repositoryHandler func(data []*models.Repository)) error {

	listOptions := &github.OrganizationsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	processingErrors := make([]error, 0)
	for {
		organizations, response, err := client.Organizations.ListAll(context.Background(), listOptions)
		if err != nil {
			processingErrors = append(processingErrors, err)
		}

		for _, item := range organizations {
			err := c.processRepositoriesInOrganization(client, item, repositoryHandler)
			if err != nil {
				processingErrors = append(processingErrors, err)
			}
		}

		if response.NextPage == 0 || len(organizations) == 0 {
			break
		}

		listOptions.Since = getLastOrganization(organizations)
		listOptions.Page = response.NextPage
	}

	if len(processingErrors) == 0 {
		return nil
	}
	return core.ConsolidateErrors(processingErrors)
}

func getLastOrganization(data []*github.Organization) int64 {
	lastOrganizationPosition := len(data) - 1
	return data[lastOrganizationPosition].GetID()
}

func (c *GithubSourceCodeRepositoryClient) processAllMemberOrganizations(client *github.Client,
	repositoryHandler func(data []*models.Repository)) error {
	listOptions := &github.ListOrgMembershipsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	processingErrors := make([]error, 0)
	for {
		memberOrganizations, response, err := client.Organizations.ListOrgMemberships(context.Background(), listOptions)
		if err != nil {
			processingErrors = append(processingErrors, err)
		}

		for _, item := range memberOrganizations {
			err = c.processRepositoriesInOrganization(client, item.GetOrganization(), repositoryHandler)
			if err != nil {
				processingErrors = append(processingErrors, err)
			}
		}

		if response.NextPage == 0 || len(memberOrganizations) == 0 {
			break
		}

		listOptions.Page = response.NextPage
	}

	if len(processingErrors) == 0 {
		return nil
	}
	return core.ConsolidateErrors(processingErrors)
}

func (c *GithubSourceCodeRepositoryClient) processRepositoriesInOrganization(client *github.Client,
	organization *github.Organization,
	repositoryHandler func(data []*models.Repository)) error {
	opt := &github.RepositoryListByOrgOptions{
		Sort:        "full_name",
		Direction:   "asc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repositories, response, err := client.Repositories.ListByOrg(context.Background(), organization.GetLogin(), opt)
		if err != nil {
			return err
		}

		mappedData := make([]*models.Repository, 0)
		for _, item := range repositories {
			mappedItem := mapRepository(c.host, organization, item)
			mappedData = append(mappedData, mappedItem)
		}

		if repositoryHandler != nil {
			repositoryHandler(mappedData)
		}

		if response.NextPage == 0 {
			break
		}

		opt.Page = response.NextPage
	}

	return nil
}
