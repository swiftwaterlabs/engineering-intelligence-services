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
)

func ExtractRepositories(host string,
	options *models.RepositoryProcessingOptions,
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
			options *models.RepositoryProcessingOptions,
			configurationService configuration.ConfigurationService,
			dataHub messaging.MessageHub) {
			defer hostWaitGroup.Done()

			log.Printf("Procesing host:%s", host.Id)
			err := processHostRepositories(host, options, configurationService, dataHub)
			if err != nil {
				log.Printf("Error when processing host:%s|%s", host.Id, err)
			}
		}

		go processor(item, options, configurationService, dataHub)
	}
	hostWaitGroup.Wait()

	return nil
}

func processHostRepositories(host *models.Host,
	options *models.RepositoryProcessingOptions,
	configurationService configuration.ConfigurationService,
	dataHub messaging.MessageHub) error {

	client, err := clients.NewSourceCodeRepositoryClient(host)
	if err != nil {
		return err
	}

	publishingQueue := configurationService.GetValue("engineering_intelligence_prd_ingestion_queue")

	handlers := getDataHandlers(host, publishingQueue, dataHub)
	log.Printf("Sending repositories from %s to %s", host.Name, publishingQueue)
	processingErr := client.ProcessRepositories(configurationService, options, handlers)

	return processingErr
}

func getDataHandlers(host *models.Host, publishingQueue string, dataHub messaging.MessageHub) *clients.RepositoryHandlers {
	handlers := &clients.RepositoryHandlers{}
	repoCounter := 0
	handlers.Repository = func(data []*models.Repository) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		repoCounter += len(data)
		log.Printf("Processed %v repositories from %s", repoCounter, host.Name)
	}

	repoOwnerCounter := 0
	handlers.Owner = func(data []*models.RepositoryOwner) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		repoOwnerCounter += len(data)
		log.Printf("Processed %v repository owners from %s", repoOwnerCounter, host.Name)
	}

	pullRequestCounter := 0
	handlers.PullRequest = func(data []*models.PullRequest) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		pullRequestCounter += len(data)
		log.Printf("Processed %v repository pull requests from %s", pullRequestCounter, host.Name)
	}

	branchRuleCounter := 0
	handlers.BranchRule = func(data []*models.BranchProtectionRule) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		branchRuleCounter += len(data)
		log.Printf("Processed %v branch rules from %s", branchRuleCounter, host.Name)
	}

	webhookCounter := 0
	handlers.Webhook = func(data []*models.Webhook) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		webhookCounter += len(data)
		log.Printf("Processed %v wehooks from %s", webhookCounter, host.Name)
	}
	return handlers
}
