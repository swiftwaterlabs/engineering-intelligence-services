package clients

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

func mapRepository(host *models.Host, organization *github.Organization, repository *github.Repository) *models.Repository {
	return &models.Repository{
		Id:                  core.MapUniqueIdentifier(host.Id, organization.GetLogin(), repository.GetName()),
		Type:                "repository",
		Organization:        mapOrganization(host, organization),
		Name:                repository.GetName(),
		DefaultBranch:       repository.GetDefaultBranch(),
		Url:                 repository.GetHTMLURL(),
		Visibility:          repository.GetVisibility(),
		CreatedAt:           repository.GetCreatedAt().Time,
		UpdatedAt:           repository.GetUpdatedAt().Time,
		ContentsLastUpdated: repository.GetPushedAt().Time,
		IsForkedRepository:  repository.GetFork(),
		ForksCount:          repository.GetForksCount(),
		RawData:             repository,
	}
}

func mapOrganization(host *models.Host, organization *github.Organization) models.Organization {
	return models.Organization{
		Id:       core.MapUniqueIdentifier(host.Id, organization.GetLogin()),
		Type:     "organization",
		Host:     host.Id,
		HostType: host.SubType,
		Url:      organization.GetHTMLURL(),
		Name:     organization.GetLogin(),
		RawData:  organization,
	}
}

func mapRepositoryOwner(repository *models.Repository, pattern string, owner string, parentOwner string) *models.RepositoryOwner {
	ownerData := &models.RepositoryOwner{
		Id:             core.MapUniqueIdentifier(repository.Organization.Host, repository.Organization.Name, repository.Name, parentOwner, owner, pattern),
		Type:           "repository-owner",
		Organization:   repository.Organization,
		RepositoryName: repository.Name,
		Pattern:        pattern,
		Owner:          owner,
		ParentOwner:    parentOwner,
	}
	return ownerData
}

func mapPullRequest(repository *models.Repository,
	pullRequest *github.PullRequest,
	reviews []*github.PullRequestReview) *models.PullRequest {
	ownerData := &models.PullRequest{
		Id:              core.MapUniqueIdentifier(repository.Organization.Host, repository.Organization.Name, repository.Name, fmt.Sprint(pullRequest.GetNumber())),
		Type:            "pull-request",
		Repository:      repository,
		TargetBranch:    pullRequest.GetBase().GetRef(),
		Url:             pullRequest.GetHTMLURL(),
		Title:           pullRequest.GetTitle(),
		Status:          pullRequest.GetTitle(),
		IsMerged:        pullRequest.GetMerged(),
		ClosedAt:        pullRequest.GetClosedAt(),
		CreatedAt:       pullRequest.GetCreatedAt(),
		CreatedBy:       pullRequest.GetUser().GetLogin(),
		Reviews:         mapPullRequestReview(reviews),
		RawData:         pullRequest,
		RawReviewerData: reviews,
	}
	return ownerData
}

func mapPullRequestReview(reviews []*github.PullRequestReview) []*models.PullRequestReview {
	result := make([]*models.PullRequestReview, 0)

	for _, item := range reviews {
		mappedReview := &models.PullRequestReview{
			Reviewer:   item.GetUser().GetLogin(),
			Status:     item.GetState(),
			ReviewedAt: item.GetSubmittedAt(),
		}

		result = append(result, mappedReview)
	}

	return result
}
