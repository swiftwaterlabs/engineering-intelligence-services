package clients

import (
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
		RawData:             repository,
	}
}

func mapOrganization(host *models.Host, organization *github.Organization) models.Organization {
	return models.Organization{
		Id:          core.MapUniqueIdentifier(host.Id, organization.GetLogin()),
		Type:        "organization",
		Host:        host.Id,
		HostType:    host.SubType,
		Url:         organization.GetHTMLURL(),
		Name:        organization.GetLogin(),
		Description: organization.GetDescription(),
		RawData:     organization,
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
