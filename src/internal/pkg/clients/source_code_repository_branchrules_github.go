package clients

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
)

func (c *GithubSourceCodeRepositoryClient) processBranchProtectionrulesForRepository(client *github.Client,
	repository *models.Repository) []*models.BranchProtectionRule {

	rule, response, err := client.Repositories.GetBranchProtection(context.Background(), repository.Organization.Name, repository.Name, repository.DefaultBranch)
	if err != nil && response.StatusCode != 404 {
		log.Printf("Unable to obtain branch protection rules for %s", repository.Url)
		return make([]*models.BranchProtectionRule, 0)
	}
	if rule == nil {
		return make([]*models.BranchProtectionRule, 0)
	}
	mappedRule := mapBranchProtectionRule(repository, repository.DefaultBranch, rule)
	return []*models.BranchProtectionRule{mappedRule}
}
