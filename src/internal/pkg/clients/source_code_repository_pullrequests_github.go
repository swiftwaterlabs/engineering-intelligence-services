package clients

import (
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"time"
)

func (c *GithubSourceCodeRepositoryClient) processPullRequestsForRepository(client *github.Client, item *models.Repository, since *time.Time) []*models.PullRequest {
	return make([]*models.PullRequest, 0)
}
