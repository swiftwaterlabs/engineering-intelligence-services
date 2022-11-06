package orchestration

import (
	"context"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/clients"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"log"
	"sync"
	"time"
)

func ExtractRepositories(host string,
	since *time.Time,
	configurationService configuration.ConfigurationService,
	hostRepository repositories.HostRepository,
	dataHub messaging.MessageHub) error {

	hosts, err := getHosts(host, models.HostTypeSourceCodeRepository, hostRepository)
	if err != nil {
		return err
	}

	hostWaitGroup := sync.WaitGroup{}
	for _, item := range hosts {
		hostWaitGroup.Add(1)
		processor := func(host *models.Host,
			since *time.Time,
			configurationService configuration.ConfigurationService,
			dataHub messaging.MessageHub) {
			defer hostWaitGroup.Done()

			log.Printf("Procesing host:%s", host.Id)
			err := processHostRepositories(host, since, configurationService, dataHub)
			if err != nil {
				log.Printf("Error when processing host:%s|%s", host.Id, err)
			}
		}

		go processor(item, since, configurationService, dataHub)
	}
	hostWaitGroup.Wait()

	return nil
}

func processHostRepositories(host *models.Host,
	since *time.Time,
	configurationService configuration.ConfigurationService,
	dataHub messaging.MessageHub) error {

	hostSecret := configurationService.GetSecret(host.ClientSecretName)
	client, err := clients.GetGitHubClient(host.SubType, host.BaseUrl, host.AuthenticationType, hostSecret)
	if err != nil {
		return err
	}

	user, response, err := client.Users.Get(context.Background(), "jrolstad")
	if err != nil {
		return err
	}
	log.Printf("Users:%s|Response Code:%v", user.GetURL(), response.StatusCode)
	return nil
}
