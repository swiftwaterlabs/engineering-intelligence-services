package orchestration

import (
	"github.com/google/uuid"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
)

func ProcessWebhookEvent(headers map[string]string, event string, configurationService configuration.ConfigurationService, dataHub messaging.MessageHub) error {
	var eventAsInterface interface{}
	eventAsInterface = core.MapFromJson(event, eventAsInterface)
	eventData := &models.WebhookEvent{
		Id:      getWebhookUniqueIdentifier(headers),
		Type:    "webhook-event",
		Source:  getWebhookSource(headers),
		Headers: headers,
		RawData: eventAsInterface,
	}

	publishingQueue := configurationService.GetValue("engineering_intelligence_prd_ingestion_queue")
	toPublish := []*models.WebhookEvent{
		eventData,
	}
	err := dataHub.SendBulk(core.ToInterfaceSlice(toPublish), publishingQueue)

	return err
}

func getWebhookUniqueIdentifier(headers map[string]string) string {
	if headers["X-GitHub-Delivery"] != "" {
		return headers["X-GitHub-Delivery"]
	}
	
	return uuid.New().String()
}

func getWebhookSource(headers map[string]string) string {
	if headers["X-GitHub-Event"] != "" {
		return "github"
	}

	if headers["X-SonarQube-Project"] != "" {
		return "sonarqube"
	}

	return ""
}
