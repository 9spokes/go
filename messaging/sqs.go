package messaging

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQS is an SQS structure
type SQS struct {
	SQS *sqs.SQS
	URL string
}

// Connect is an SQS connection convenience function
func (_sqs *SQS) Connect(url string) error {

	_sqs.SQS = sqs.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	_sqs.URL = url
	fmt.Printf("Queue URL: %s\n", url)
	return nil
}

// SendMessage is an SQS convenience method to send a message to a given queue name
func (_sqs *SQS) SendMessage(queue string, message Message) error {

	_, err := _sqs.SQS.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(message.Body)),
		QueueUrl:     &queue,
	})

	return err
}

// ReceiveMessages is an SQS convenience method to retrieve messages from a queue
func (_sqs *SQS) ReceiveMessages(queue string) (<-chan Message, error) {

	ret := make(chan Message)
	go func() {

		for {

			output, err := _sqs.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				},
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            &queue,
				MaxNumberOfMessages: aws.Int64(1),
				VisibilityTimeout:   aws.Int64(20), // 20 seconds
				WaitTimeSeconds:     aws.Int64(20),
			})
			if err != nil {
				panic(err)
			}

			for _, message := range output.Messages {
				fmt.Printf("Publishing message\n")
				ret <- Message{*message.ReceiptHandle, *message.MessageId, []byte(*message.Body)}
				fmt.Printf("Deleting message %s from queue %s\n", *message.ReceiptHandle, _sqs.URL)
				err := _sqs.DeleteMessage(*message.ReceiptHandle)
				if err != nil {
					panic(fmt.Sprintf("Error deleting message: %s\n", err.Error()))
				}
			}
		}
	}()

	return ret, nil
}

// DeleteMessage removes a message from an SQS Queue
func (_sqs *SQS) DeleteMessage(id string) error {

	if id == "" {
		return errors.New("Message ID is required")
	}

	_, err := _sqs.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &_sqs.URL,
		ReceiptHandle: &id,
	})

	return err
}
