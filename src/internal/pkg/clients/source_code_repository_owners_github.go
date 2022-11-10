package clients

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
	"strings"
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

func (c *GithubSourceCodeRepositoryClient) resolveRepositoryOwners(client *github.Client, repository *models.Repository,
	codeOwners map[string]map[string]*codeOwnerData) []*models.RepositoryOwner {

	repositoryCodeOwner := c.coalesceCodeOwners(codeOwners[repository.Name]["CODEOWNERS"],
		codeOwners[repository.Name]["docs/CODEOWNERS"],
		codeOwners[repository.Name][".github/CODEOWNERS"])
	organizationCodeOwner := c.coalesceCodeOwners(
		codeOwners["sfdc-codeowners"][fmt.Sprintf("%s/CODEOWNERS", repository.Name)],
		codeOwners["sfdc-codeowners"]["sfdc-codeowners-uo/CODEOWNERS"])

	repositoryCodeOwners := make([]*models.RepositoryOwner, 0)
	if repositoryCodeOwner != nil {
		data := c.parseCodeOwners(repository, repositoryCodeOwner.Contents)
		repositoryCodeOwners = append(repositoryCodeOwners, data...)
	}

	organizationCodeOwners := make([]*models.RepositoryOwner, 0)
	if organizationCodeOwner != nil {
		data := c.parseCodeOwners(repository, organizationCodeOwner.Contents)
		organizationCodeOwners = append(organizationCodeOwners, data...)
	}

	c.applyOrganizationDefaults(repositoryCodeOwners, organizationCodeOwners)

	if len(repositoryCodeOwners) > 0 {
		return repositoryCodeOwners
	}

	return organizationCodeOwners
}

func (c *GithubSourceCodeRepositoryClient) applyOrganizationDefaults(repositoryCodeOwners []*models.RepositoryOwner,
	organizationCodeOwners []*models.RepositoryOwner) {
	for _, item := range repositoryCodeOwners {
		for _, orgItem := range organizationCodeOwners {
			if item.ParentOwner == "" {
				item.ParentOwner = orgItem.ParentOwner
			}
		}
	}
}

func (c *GithubSourceCodeRepositoryClient) coalesceCodeOwners(items ...*codeOwnerData) *codeOwnerData {
	for _, value := range items {
		if value != nil {
			return value
		}
	}

	return nil
}

func (c *GithubSourceCodeRepositoryClient) parseCodeOwners(repository *models.Repository, contents string) []*models.RepositoryOwner {
	if strings.TrimSpace(contents) == "" {
		return make([]*models.RepositoryOwner, 0)
	}

	owners := make(map[string][]*models.RepositoryOwner, 0)

	linesInFile := strings.Split(contents, "\n")

	const parentOwnerLinePrefix = "#GUSINFO:"
	const commentPrefix = "#"

	parentOwner := ""
	for _, line := range linesInFile {
		cleanLine := strings.TrimSpace(line)
		if strings.HasPrefix(cleanLine, parentOwnerLinePrefix) {
			delimitedValues := strings.TrimSpace(strings.ReplaceAll(cleanLine, parentOwnerLinePrefix, ""))
			splitValues := strings.Split(delimitedValues, ",")

			parentOwner = core.GetValueAt(splitValues, 0)
		}

		if owners[parentOwner] == nil {
			owners[parentOwner] = make([]*models.RepositoryOwner, 0)
		}

		if strings.HasPrefix(cleanLine, commentPrefix) || cleanLine == "" {
			continue
		}

		ownerParts := strings.Fields(cleanLine)
		pattern := core.GetValueAt(ownerParts, 0)

		patternOwners := ownerParts[1:]
		if len(patternOwners) == 0 {
			ownerData := mapRepositoryOwner(repository, pattern, "", parentOwner)
			owners[parentOwner] = append(owners[parentOwner], ownerData)
		} else {
			for _, item := range patternOwners {
				ownerData := mapRepositoryOwner(repository, pattern, item, parentOwner)
				owners[parentOwner] = append(owners[parentOwner], ownerData)
			}
		}
	}
	return c.mapRepositoryOwnersToSlice(owners)

}

func (c *GithubSourceCodeRepositoryClient) mapRepositoryOwnersToSlice(data map[string][]*models.RepositoryOwner) []*models.RepositoryOwner {
	results := make([]*models.RepositoryOwner, 0)

	for _, value := range data {
		results = append(results, value...)
	}

	return results
}
