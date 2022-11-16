package clients

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
	"time"
)

func (c *GithubSourceCodeRepositoryClient) processPullRequestsForRepository(client *github.Client,
	repository *models.Repository,
	since *time.Time) []*models.PullRequest {
	options := &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	result := make([]*models.PullRequest, 0)

	for {
		pullRequests, response, err := client.PullRequests.List(context.Background(), repository.Organization.Name, repository.Name, options)
		if err != nil {
			log.Printf("Unable to retrieve pull requests for %s", repository.Url)
		}

		for _, item := range pullRequests {

			if since != nil && item.GetUpdatedAt().Before(*since) {
				break
			}

			reviews := c.getPullRequestReviews(client, repository, item.GetNumber(), item.GetURL())
			files := c.getPullRequestFiles(client, repository, item.GetNumber(), item.GetURL())
			mappedPullRequest := mapPullRequest(repository, item, reviews, files)
			result = append(result, mappedPullRequest)

		}

		if response == nil || response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}
	return result
}

func (s *GithubSourceCodeRepositoryClient) getPullRequestReviews(client *github.Client,
	repository *models.Repository,
	prNumber int,
	prUrl string) []*github.PullRequestReview {

	options := &github.ListOptions{
		Page:    0,
		PerPage: 100,
	}
	result := make([]*github.PullRequestReview, 0)

	for {
		reviews, response, err := client.PullRequests.ListReviews(context.Background(), repository.Organization.Name, repository.Name, prNumber, options)
		if err != nil {
			log.Printf("Unable to retrieve reviewers for %s", prUrl)
		}

		for _, item := range reviews {
			result = append(result, item)
		}

		if response == nil || response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	return result
}

func (s *GithubSourceCodeRepositoryClient) getPullRequestFiles(client *github.Client,
	repository *models.Repository,
	prNumber int,
	prUrl string) []*github.CommitFile {

	options := &github.ListOptions{
		Page:    0,
		PerPage: 100,
	}
	result := make([]*github.CommitFile, 0)

	for {
		files, response, err := client.PullRequests.ListFiles(context.Background(), repository.Organization.Name, repository.Name, prNumber, options)
		if err != nil {
			log.Printf("Unable to retrieve files for %s", prUrl)
		}

		for _, item := range files {
			result = append(result, item)
		}

		if response == nil || response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	return result
}
