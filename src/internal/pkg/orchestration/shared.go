package orchestration

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"strings"
)

func getHosts(host string, hostType string, hostRepository repositories.HostRepository) ([]*models.Host, error) {
	if strings.TrimSpace(host) == "" {
		return hostRepository.GetAll(hostType)
	}

	hostData, err := hostRepository.Get(host)
	if err != nil {
		return make([]*models.Host, 0), err
	}

	return []*models.Host{hostData}, nil
}
