package clients

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	sonargo "github.com/magicsong/sonargo/sonar"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

func mapRepository(host *models.Host, organization *github.Organization, repository *github.Repository) *models.Repository {
	return &models.Repository{
		Id:                  core.MapUniqueIdentifier(repository.GetURL()),
		Type:                "repository",
		Organization:        mapOrganization(host, organization),
		Name:                core.SanitizeString(repository.GetName()),
		DefaultBranch:       core.SanitizeString(repository.GetDefaultBranch()),
		Url:                 core.SanitizeString(repository.GetHTMLURL()),
		Visibility:          core.SanitizeString(repository.GetVisibility()),
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
		Id:       core.MapUniqueIdentifier(organization.GetURL()),
		Type:     "organization",
		Host:     core.SanitizeString(host.Id),
		HostType: core.SanitizeString(host.SubType),
		Url:      core.SanitizeString(organization.GetHTMLURL()),
		Name:     core.SanitizeString(organization.GetLogin()),
		RawData:  organization,
	}
}

func mapRepositoryOwner(repository *models.Repository, pattern string, owner string, parentOwner string) *models.RepositoryOwner {
	ownerData := &models.RepositoryOwner{
		Id:          core.MapUniqueIdentifier(repository.Url, pattern, owner, parentOwner),
		Type:        "repository-owner",
		Repository:  repository,
		Pattern:     core.SanitizeString(pattern),
		Owner:       core.SanitizeString(owner),
		ParentOwner: core.SanitizeString(parentOwner),
	}
	return ownerData
}

func mapPullRequest(repository *models.Repository,
	pullRequest *github.PullRequest,
	reviews []*github.PullRequestReview,
	files []*github.CommitFile) *models.PullRequest {
	ownerData := &models.PullRequest{
		Id:           core.MapUniqueIdentifier(pullRequest.GetURL()),
		Type:         "pull-request",
		Repository:   repository,
		TargetBranch: core.SanitizeString(pullRequest.GetBase().GetRef()),
		Url:          core.SanitizeString(pullRequest.GetHTMLURL()),
		Title:        core.SanitizeString(pullRequest.GetTitle()),
		Status:       core.SanitizeString(pullRequest.GetState()),
		IsMerged:     pullRequest.GetMerged() || pullRequest.MergedAt != nil || pullRequest.MergedBy != nil,
		ClosedAt:     pullRequest.GetClosedAt(),
		CreatedAt:    pullRequest.GetCreatedAt(),
		CreatedBy:    core.SanitizeString(pullRequest.GetUser().GetLogin()),
		Reviews:      mapPullRequestReview(reviews),
		Files:        mapPullRequestFiles(files),
		RawData:      pullRequest,
	}
	return ownerData
}

func mapPullRequestReview(reviews []*github.PullRequestReview) []*models.PullRequestReview {
	result := make([]*models.PullRequestReview, 0)

	for _, item := range reviews {
		mappedReview := &models.PullRequestReview{
			Reviewer:   core.SanitizeString(item.GetUser().GetLogin()),
			Status:     core.SanitizeString(item.GetState()),
			ReviewedAt: item.GetSubmittedAt(),
		}

		result = append(result, mappedReview)
	}

	return result
}

func mapPullRequestFiles(files []*github.CommitFile) []string {
	result := make([]string, 0)

	for _, item := range files {
		result = append(result, core.SanitizeString(item.GetFilename()))
	}

	return result
}

func mapBranchProtectionRule(repository *models.Repository, branchName string, rule *github.Protection) *models.BranchProtectionRule {
	return &models.BranchProtectionRule{
		Id:                           core.MapUniqueIdentifier(repository.Url, branchName),
		Type:                         "branch-rule",
		Repository:                   repository,
		Branch:                       core.SanitizeString(branchName),
		AllowForcePush:               rule.GetAllowForcePushes().Enabled,
		RequirePullRequest:           rule.RequiredPullRequestReviews != nil,
		RequirePullRequestApprovals:  rule.GetRequiredPullRequestReviews().RequireCodeOwnerReviews,
		RequiredPullRequestApprovers: rule.GetRequiredPullRequestReviews().RequiredApprovingReviewCount,
		IncludeAdministrators:        rule.GetEnforceAdmins().Enabled,
		RawData:                      rule,
	}
}

func mapWebHook(organization *models.Organization, repository *models.Repository, hook *github.Hook) *models.Webhook {
	return &models.Webhook{
		Id:           core.MapUniqueIdentifier(hook.GetURL()),
		Type:         "webhook",
		Source:       core.SanitizeString(hook.GetType()),
		Organization: organization,
		Repository:   repository,
		Events:       resolveWebhookEvents(hook),
		Target:       core.SanitizeString(resolveWebhookConfigValue(hook, "url")),
		Active:       hook.GetActive(),
		Name:         core.SanitizeString(hook.GetName()),
		RawData:      hook,
	}
}

func resolveWebhookEvents(hook *github.Hook) []string {
	if hook.Events == nil {
		return make([]string, 0)
	}

	return hook.Events
}

func resolveWebhookConfigValue(hook *github.Hook, name string) string {
	if hook.Config == nil {
		return ""
	}

	return fmt.Sprint(hook.Config[name])
}

func mapTestResult(component *sonargo.Component, data *sonargo.MeasuresSearchHistoryObject) []*models.TestResult {
	return make([]*models.TestResult, 0)
}
