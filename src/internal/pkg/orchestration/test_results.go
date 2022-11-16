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

func ExtractTestResults(host string,
	options *models.TestResultProcessingOptions,
	configurationService configuration.ConfigurationService,
	hostRepository repositories.HostRepository,
	dataHub messaging.MessageHub) error {

	hosts, err := getHosts(host, models.HostTypeAutomatedTestingPlatform, hostRepository)
	if err != nil {
		return err
	}

	hostWaitGroup := sync.WaitGroup{}
	for _, item := range hosts {
		hostWaitGroup.Add(1)
		processor := func(host *models.Host,
			options *models.TestResultProcessingOptions,
			configurationService configuration.ConfigurationService,
			dataHub messaging.MessageHub) {
			defer hostWaitGroup.Done()

			log.Printf("Procesing host:%s", host.Id)
			err := processHostTestResults(host, options, configurationService, dataHub)
			if err != nil {
				log.Printf("Error when processing host:%s|%s", host.Id, err)
			}
		}

		go processor(item, options, configurationService, dataHub)
	}
	hostWaitGroup.Wait()

	return nil
}

func processHostTestResults(host *models.Host,
	options *models.TestResultProcessingOptions,
	configurationService configuration.ConfigurationService,
	dataHub messaging.MessageHub) error {

	client, err := clients.NewTestResultClient(host)
	if err != nil {
		return err
	}

	publishingQueue := configurationService.GetValue("engineering_intelligence_prd_ingestion_queue")

	handler := getTestResultHandler(host, dataHub, publishingQueue)

	log.Printf("Sending test results from %s to %s", host.Name, publishingQueue)
	processingErr := client.ProcessTestResults(configurationService, options, handler)

	return processingErr
}

func getTestResultHandler(host *models.Host,
	dataHub messaging.MessageHub,
	publishingQueue string) func(data []*models.TestResult) {
	counter := 0
	handler := func(data []*models.TestResult) {
		toPublish := core.ToInterfaceSlice(data)
		err := dataHub.SendBulk(toPublish, publishingQueue)
		if err != nil {
			log.Println(err)
		}
		counter += len(data)
		log.Printf("Processed %v repositories from %s", counter, host.Name)
	}
	return handler
}
