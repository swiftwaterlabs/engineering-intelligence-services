package orchestration

import (
	"errors"
	"github.com/google/go-github/v48/github"
	"github.com/google/uuid"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"time"
)

func ProcessWebhookEvent(headers map[string]string,
	event string,
	authenticateEvents bool,
	eventSourceRepository repositories.EventSourceRepository,
	configurationService configuration.ConfigurationService,
	dataHub messaging.MessageHub) (bool, error) {

	eventSourceName := getWebhookSource(headers)

	if authenticateEvents {
		isEventSourceValid, err := authenticateEvent(headers, event, eventSourceName, eventSourceRepository, configurationService)
		if err != nil || !isEventSourceValid {
			return isEventSourceValid, err
		}
	}

	eventData := mapToWebhookEvent(headers, event, eventSourceName)

	signalPublishingQueue := configurationService.GetValue("engineering_intelligence_prd_ingestion_queue")
	err := dataHub.Send(eventData, signalPublishingQueue)
	if err != nil {
		return true, err
	}

	webhookEventQueue := configurationService.GetValue("engineering_intelligence_prd_webhook_event_queue")
	err = dataHub.Send(eventData, webhookEventQueue)

	return true, err
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

func authenticateEvent(headers map[string]string,
	event string,
	sourceName string,
	eventSourceRepository repositories.EventSourceRepository,
	configurationService configuration.ConfigurationService) (bool, error) {

	actualHash := ""
	if sourceName == "github" {
		actualHash = headers["X-Hub-Signature-256"]
	} else if sourceName == "sonarqube" {
		actualHash = headers["X-Sonar-Webhook-HMAC-SHA256"]
	} else {
		return false, errors.New("unrecognized source")
	}

	bodyAsBytes := []byte(event)

	eventSources, err := eventSourceRepository.GetAllActive()
	if err != nil {
		return false, err
	}

	for _, source := range eventSources {
		for _, secret := range source.Secrets {
			if secret.Active {
				secretValue := configurationService.GetSecret(secret.SecretName)
				secretValueAsBytes := []byte(secretValue)

				// Use the go-github implementation since they've already done the hard work
				err := github.ValidateSignature(actualHash, bodyAsBytes, secretValueAsBytes)
				if err == nil {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func mapToWebhookEvent(headers map[string]string, event string, eventSourceName string) *models.WebhookEvent {
	var webhookEvent interface{}
	core.MapFromJson(event, &webhookEvent)

	eventData := &models.WebhookEvent{
		Id:         getWebhookUniqueIdentifier(headers),
		Type:       "webhook-event",
		Source:     eventSourceName,
		EventType:  getWebhookEventType(headers),
		ReceivedAt: time.Now(),
		Headers:    headers,
		RawData:    webhookEvent,
	}
	return eventData
}

func getWebhookUniqueIdentifier(headers map[string]string) string {
	if headers["X-GitHub-Delivery"] != "" {
		return headers["X-GitHub-Delivery"]
	}

	return uuid.New().String()
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
