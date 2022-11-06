package orchestration

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"time"
)

func ExtractRepositories(host string,
	since *time.Time,
	configurationService configuration.ConfigurationService,
	hostRepository repositories.HostRepository,
	dataHub messaging.MessageHub) error {
	return nil
}
