package orchestration

import (
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/clients"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"log"
	"sync"
	"time"
)

func ExtractRepositories(host string,
	since *time.Time,
	includeRepositoryDetails bool,
	includeOwners bool,
	includePullRequests bool,
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
			err := processHostRepositories(host, includeRepositoryDetails, includeOwners, includePullRequests, since, configurationService, dataHub)
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
	includeRepositoryDetails bool,
	includeOwners bool,
	includePullRequests bool,
	since *time.Time,
	configurationService configuration.ConfigurationService,
	dataHub messaging.MessageHub) error {

	client, err := clients.NewSourceCodeRepositoryClient(host)
	if err != nil {
		return err
	}

	publishingQueue := configurationService.GetValue("engineering_intelligence_prd_ingestion_queue")

	log.Printf("Sending repositories from %s to %s", host.Name, publishingQueue)
	repoCounter := 0
	repositoryHandler := func(data []*models.Repository) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		repoCounter += len(data)
		log.Printf("Processed %v repositories from %s", repoCounter, host.Name)
	}

	repoOwnerCounter := 0
	ownerHandler := func(data []*models.RepositoryOwner) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		repoOwnerCounter += len(data)
		log.Printf("Processed %v repository owners from %s", repoOwnerCounter, host.Name)
	}

	pullRequestCounter := 0
	pullRequestHandler := func(data []*models.PullRequest) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		pullRequestCounter += len(data)
		log.Printf("Processed %v repository pull requests from %s", pullRequestCounter, host.Name)
	}

	processingErr := client.ProcessRepositories(configurationService, includeRepositoryDetails, includeOwners, includePullRequests, since, repositoryHandler, ownerHandler, pullRequestHandler)

	return processingErr
}
