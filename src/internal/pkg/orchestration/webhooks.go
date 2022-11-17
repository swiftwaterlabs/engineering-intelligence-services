package orchestration

import (
	"github.com/google/uuid"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"time"
)

func ProcessWebhookEvent(headers map[string]string, event string, configurationService configuration.ConfigurationService, dataHub messaging.MessageHub) error {
	var webhookEvent interface{}
	core.MapFromJson(event, &webhookEvent)
	eventData := &models.WebhookEvent{
		Id:         getWebhookUniqueIdentifier(headers),
		Type:       "webhook-event",
		Source:     getWebhookSource(headers),
		EventType:  getWebhookEventType(headers),
		ReceivedAt: time.Now(),
		Headers:    headers,
		RawData:    webhookEvent,
	}

	publishingQueue := configurationService.GetValue("engineering_intelligence_prd_ingestion_queue")
	err := dataHub.Send(eventData, publishingQueue)

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

func getWebhookEventType(headers map[string]string) string {
	if headers["X-GitHub-Event"] != "" {
		return headers["X-GitHub-Event"]
	}

	if headers["X-SonarQube-Project"] != "" {
		return "analysis_complete"
	}

	return ""
}
