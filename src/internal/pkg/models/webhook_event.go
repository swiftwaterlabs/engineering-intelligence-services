package models

import "time"

type WebhookEvent struct {
	Id         string
	Type       string
	Source     string
	EventType  string
	ReceivedAt time.Time
	Headers    map[string]string
	RawData    interface{}
}
