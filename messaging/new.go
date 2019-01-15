package messaging

import "fmt"

// Transport is a messaging protocol transport
type Transport interface {
	Connect(string) error
	SendMessage(string, Message) error
	DeleteMessage(string) error
	CreateQueue(string, map[string]interface{}) error
	ReceiveMessages(string) (<-chan Message, error)
}

// Message is an abstract message structure
type Message struct {
	ID            string
	CorrelationID string
	Body          []byte
}

// New is a function that creates a new transport type
func New(transport string) (Transport, error) {

	if transport == "amqp" {
		return &AMQP{}, nil
	}

	// Commented out to reduce package sizes
	// if transport == "sqs" {
	// 	return &SQS{}, nil
	// }
	return nil, fmt.Errorf("Unknown transport type '%s'", transport)
}
