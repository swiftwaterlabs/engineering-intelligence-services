package clients

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

type SourceCodeRepositoryOwnerResolver interface {
	ResolveRepositoryOwners(repository *models.Repository, codeOwners map[string]map[string]*codeOwnerData) []*models.RepositoryOwner
}

func NewSourceCodeRepositoryOwnerResolver(host *models.Host) SourceCodeRepositoryOwnerResolver {
	if host.Options["RepositoryOwnerResolver"] == "sfdc" {
		return &SfdcSourceCodeRepositoryOwnerResolver{}
	}
	return &DefaultSourceCodeRepositoryOwnerResolver{}
}

func coalesceCodeOwners(items ...*codeOwnerData) *codeOwnerData {
	for _, value := range items {
		if value != nil {
			return value
		}
	}

	return nil
}
