package messaging

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
)

type MessageHub interface {
	Send(toSend interface{}, target string) error
	SendBulk(toSend []interface{}, target string) error
	Receive(target string, handler func(message interface{})) error
}

func NewMessageHub(config *configuration.AppConfig) MessageHub {
	hub := new(SqsMessageHub)

	session := configuration.GetAwsSession(config)
	hub.sqs = sqs.New(session)

	return hub
}
