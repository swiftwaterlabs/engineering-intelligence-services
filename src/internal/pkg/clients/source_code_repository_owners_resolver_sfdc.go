package clients

import (
	"fmt"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"strings"
)

type SfdcSourceCodeRepositoryOwnerResolver struct {
}

func (c *SfdcSourceCodeRepositoryOwnerResolver) ResolveRepositoryOwners(repository *models.Repository,
	codeOwners map[string]map[string]*codeOwnerData) []*models.RepositoryOwner {

	repositoryCodeOwner := coalesceCodeOwners(codeOwners[repository.Name]["CODEOWNERS"],
		codeOwners[repository.Name]["docs/CODEOWNERS"],
		codeOwners[repository.Name][".github/CODEOWNERS"])
	organizationCodeOwner := coalesceCodeOwners(
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

func (c *SfdcSourceCodeRepositoryOwnerResolver) applyOrganizationDefaults(repositoryCodeOwners []*models.RepositoryOwner,
	organizationCodeOwners []*models.RepositoryOwner) {
	for _, item := range repositoryCodeOwners {
		for _, orgItem := range organizationCodeOwners {
			if item.ParentOwner == "" {
				item.ParentOwner = orgItem.ParentOwner
			}
		}
	}
}

func (c *SfdcSourceCodeRepositoryOwnerResolver) parseCodeOwners(repository *models.Repository, contents string) []*models.RepositoryOwner {
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

func (c *SfdcSourceCodeRepositoryOwnerResolver) mapRepositoryOwnersToSlice(data map[string][]*models.RepositoryOwner) []*models.RepositoryOwner {
	results := make([]*models.RepositoryOwner, 0)

	for _, value := range data {
		results = append(results, value...)
	}

	return results
}
