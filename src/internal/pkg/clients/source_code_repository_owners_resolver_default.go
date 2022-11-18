package clients

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"strings"
)

type DefaultSourceCodeRepositoryOwnerResolver struct {
}

func (c *DefaultSourceCodeRepositoryOwnerResolver) ResolveRepositoryOwners(repository *models.Repository,
	codeOwners map[string]map[string]*codeOwnerData) []*models.RepositoryOwner {
	repositoryCodeOwner := coalesceCodeOwners(codeOwners[repository.Name]["CODEOWNERS"],
		codeOwners[repository.Name]["docs/CODEOWNERS"],
		codeOwners[repository.Name][".github/CODEOWNERS"])

	if repositoryCodeOwner == nil {
		return make([]*models.RepositoryOwner, 0)
	}

	result := c.parseCodeOwners(repository, repositoryCodeOwner.Contents)

	return result
}

func (c *DefaultSourceCodeRepositoryOwnerResolver) parseCodeOwners(repository *models.Repository, contents string) []*models.RepositoryOwner {
	if strings.TrimSpace(contents) == "" {
		return make([]*models.RepositoryOwner, 0)
	}

	linesInFile := strings.Split(contents, "\n")
	const commentPrefix = "#"

	results := make([]*models.RepositoryOwner, 0)
	for _, line := range linesInFile {
		cleanLine := strings.TrimSpace(line)

		if strings.HasPrefix(cleanLine, commentPrefix) || cleanLine == "" {
			continue
		}

		ownerParts := strings.Fields(cleanLine)
		pattern := core.GetValueAt(ownerParts, 0)

		patternOwners := ownerParts[1:]
		if len(patternOwners) == 0 {
			ownerData := mapRepositoryOwner(repository, pattern, "", "")
			results = append(results, ownerData)
		} else {
			for _, item := range patternOwners {
				ownerData := mapRepositoryOwner(repository, pattern, item, "")
				results = append(results, ownerData)
			}
		}
	}

	return results
}
