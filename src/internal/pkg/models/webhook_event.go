package models

type WebhookEvent struct {
	Id      string
	Type    string
	Source  string
	Headers map[string]string
	RawData interface{}
}
