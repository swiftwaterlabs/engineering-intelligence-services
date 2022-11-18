package orchestration

import (
	"errors"
	"github.com/google/go-github/v48/github"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/messaging"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/repositories"
	"log"
	"strings"
)

func ListenForWebhookEvents(configurationService configuration.ConfigurationService,
	hostRepository repositories.HostRepository,
	dataHub messaging.MessageHub) error {
	webhookEventQueue := configurationService.GetValue("engineering_intelligence_prd_webhook_event_queue")

	handler := func(item interface{}) {
		data := core.MapToJson(item)
		eventData := &models.WebhookEvent{}
		core.MapFromJson(data, eventData)

		err := EnhanceWebhookEvent(eventData, configurationService, hostRepository, dataHub)
		if err != nil {
			log.Printf("error when processing webhook event|%s", err)
		}
	}

	err := dataHub.Receive(webhookEventQueue, handler)

	return err
}

func EnhanceWebhookEvent(item *models.WebhookEvent,
	configurationService configuration.ConfigurationService,
	hostRepository repositories.HostRepository,
	dataHub messaging.MessageHub) error {

	if item == nil {
		return errors.New("unable to process empty event")
	}

	if isGithubPullRequestEvent(item) {
		return processGithubPullRequestEvent(item)
	}

	if isBranchProtectionRuleEvent(item) {
		return processBranchProtectionRuleEvent(item)
	}
	return nil
}

func isGithubPullRequestEvent(item *models.WebhookEvent) bool {
	return strings.EqualFold(item.Source, models.WebhookEventSourceGithub) && strings.EqualFold(item.EventType, "pull_request")
}

func processGithubPullRequestEvent(item *models.WebhookEvent) error {
	webhook := &github.PullRequestEvent{}
	webhookData := core.MapToJson(item.RawData)
	core.MapFromJson(webhookData, webhook)

	if strings.EqualFold(webhook.GetAction(), "closed") {

	}

	return nil
}

func isBranchProtectionRuleEvent(item *models.WebhookEvent) bool {
	return strings.EqualFold(item.Source, models.WebhookEventSourceGithub) && strings.EqualFold(item.EventType, "branch_protection_rule")
}

func processBranchProtectionRuleEvent(item *models.WebhookEvent) error {
	webhook := &github.BranchProtectionRuleEvent{}
	webhookData := core.MapToJson(item.RawData)
	core.MapFromJson(webhookData, webhook)

	return nil
}
