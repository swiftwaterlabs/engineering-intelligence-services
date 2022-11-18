package clients

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"log"
)

func (s *GithubSourceCodeRepositoryClient) getCodeOwnersForOrganization(client *github.Client,
	organization *github.Organization) (map[string]map[string]*codeOwnerData, error) {
	searchOptions := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	results := make(map[string]map[string]*codeOwnerData, 0)
	query := fmt.Sprintf("filename:CODEOWNERS org:%s", organization.GetLogin())
	for {
		result, response, err := client.Search.Code(context.Background(), query, searchOptions)
		if err != nil {
			return results, err
		}
		if result.GetIncompleteResults() {
			log.Printf("Incomplete search results for %s on page %v", query, searchOptions.Page)
		}

		for _, item := range result.CodeResults {
			repositoryName := item.GetRepository().GetName()
			path := item.GetPath()

			if results[repositoryName] == nil {
				results[repositoryName] = make(map[string]*codeOwnerData, 0)
			}

			data := &codeOwnerData{
				Organization: organization.GetLogin(),
				Repository:   repositoryName,
				Path:         path,
			}
			results[repositoryName][data.Path] = data
		}

		if response.NextPage == 0 {
			break
		}
		searchOptions.Page = response.NextPage
	}

	s.getCodeOwnersContent(client, results)
	return results, nil
}

type codeOwnerData struct {
	Organization string
	Repository   string
	Path         string
	Contents     string
}

func (s *GithubSourceCodeRepositoryClient) getCodeOwnersContent(client *github.Client,
	organizationCodeOwners map[string]map[string]*codeOwnerData) {
	options := &github.RepositoryContentGetOptions{}
	for _, repositoryCodeOwners := range organizationCodeOwners {
		for _, file := range repositoryCodeOwners {
			fileContent, _, _, err := client.Repositories.GetContents(context.Background(), file.Organization, file.Repository, file.Path, options)
			if err == nil && fileContent != nil {
				content, contentErr := fileContent.GetContent()
				if contentErr == nil {
					file.Contents = content
				}
			}
		}
	}
}
