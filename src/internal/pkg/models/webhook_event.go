package models

import "time"

const WebhookEventSourceGithub = "github"
const WebhookEventSourceSonarqube = "sonarqube"

type WebhookEvent struct {
	Id         string
	Type       string
	Source     string
	EventType  string
	ReceivedAt time.Time
	Headers    map[string]string
	RawData    interface{}
}
