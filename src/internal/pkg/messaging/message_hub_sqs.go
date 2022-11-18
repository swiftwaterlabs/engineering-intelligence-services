package messaging

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/core"
)

type SqsMessageHub struct {
	sqs *sqs.SQS
}

func (hub *SqsMessageHub) Send(toSend interface{}, target string) error {
	message, mapError := hub.mapToSqsSendMessage(hub.sqs, toSend, target)
	if mapError != nil {
		return mapError
	}

	_, sendError := hub.sqs.SendMessage(message)
	if sendError != nil {
		return sendError
	}

	return nil
}

func (hub *SqsMessageHub) SendBulk(toSend []interface{}, target string) error {
	const maxBatchSizeInBytes = 262144 //256KB
	const maxBatchSizeInItems = 10
	batches, segmentErr := core.SegmentByJsonByteLength(toSend, maxBatchSizeInBytes, maxBatchSizeInItems)
	if segmentErr != nil {
		return segmentErr
	}

	sendErrors := make([]error, 0)
	for _, batch := range batches {
		message, mapError := hub.mapToSqsSendMessageBatch(hub.sqs, batch, target)
		if mapError != nil {
			sendErrors = append(sendErrors, mapError)
			continue
		}

		output, sendError := hub.sqs.SendMessageBatch(message)
		if sendError != nil {
			sendErrors = append(sendErrors, sendError)
			continue
		}
		if len(output.Failed) > 0 {
			failedError := errors.New("at least one item in the batch failed to send")
			sendErrors = append(sendErrors, failedError)
		}
	}

	if len(sendErrors) > 0 {
		return core.ConsolidateErrors(sendErrors)
	}
	return nil
}

func (hub *SqsMessageHub) Receive(target string, handler func(message interface{})) error {
	urlResult, queueUrlErr := hub.sqs.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(target)})
	if queueUrlErr != nil {
		return queueUrlErr
	}

	options := &sqs.ReceiveMessageInput{
		QueueUrl: urlResult.QueueUrl,
	}

	for {
		result, err := hub.sqs.ReceiveMessage(options)
		if err != nil {
			break
		}

		for _, item := range result.Messages {
			var data interface{}
			core.MapFromJson(aws.StringValue(item.Body), &data)
			handler(data)

			deleteOptions := &sqs.DeleteMessageInput{
				QueueUrl:      urlResult.QueueUrl,
				ReceiptHandle: item.ReceiptHandle,
			}
			hub.sqs.DeleteMessage(deleteOptions)
		}
	}

	return nil
}
func (hub *SqsMessageHub) mapToSqsSendMessage(sqsInstance *sqs.SQS, toMap interface{}, queueName string) (*sqs.SendMessageInput, error) {
	urlResult, queueUrlErr := sqsInstance.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
	if queueUrlErr != nil {
		return nil, queueUrlErr
	}

	input := new(sqs.SendMessageInput)
	input.MessageBody = aws.String(core.MapToJson(toMap))
	input.QueueUrl = urlResult.QueueUrl

	return input, nil
}

func (hub *SqsMessageHub) mapToSqsSendMessageBatch(sqsInstance *sqs.SQS, toMap []interface{}, queueName string) (*sqs.SendMessageBatchInput, error) {
	urlResult, queueUrlErr := sqsInstance.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
	if queueUrlErr != nil {
		return nil, queueUrlErr
	}

	input := &sqs.SendMessageBatchInput{
		Entries:  hub.mapToSqsSendMessageBatchRequestEntry(toMap),
		QueueUrl: urlResult.QueueUrl,
	}

	return input, nil
}

func (hub *SqsMessageHub) mapToSqsSendMessageBatchRequestEntry(toMap []interface{}) []*sqs.SendMessageBatchRequestEntry {
	entries := make([]*sqs.SendMessageBatchRequestEntry, len(toMap))
	for index, item := range toMap {
		entries[index] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(uuid.New().String()),
			MessageBody: aws.String(core.MapToJson(item)),
		}
	}

	return entries

}
