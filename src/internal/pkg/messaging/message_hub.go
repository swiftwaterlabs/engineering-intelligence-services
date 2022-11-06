package messaging

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
)

type MessageHub interface {
	Send(toSend interface{}, target string) error
	SendBulk(toSend []interface{}, target string) error
}

func NewMessageHub(config *configuration.AppConfig) MessageHub {
	hub := new(SqsMessageHub)

	session := configuration.GetAwsSession(config)
	hub.sqs = sqs.New(session)

	return hub
}
