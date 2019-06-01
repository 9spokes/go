package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

// AMQP is an AMQP structure
type AMQP struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// Connect is an AMQP connection convenience function
func (_amqp *AMQP) Connect(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	_amqp.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	_amqp.Channel = ch

	return nil
}

// SendMessage is an AMQP convenience method to send a message to a given queue name
func (_amqp *AMQP) SendMessage(queue string, message Message) error {

	err := _amqp.Channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          message.Body,
			CorrelationId: message.CorrelationID,
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to send message: %s", err.Error())
	}

	return nil
}

// DeleteMessage is an AMQP convenience method which does nothing, as AMQP does not support message deletion
func (_amqp *AMQP) DeleteMessage(id string) error {
	// No body because AMQP does not support message deletion without consumption
	return nil

}

// CreateQueue creates a new message with the given name and attributes
func (_amqp *AMQP) CreateQueue(name string, attributes map[string]interface{}) error {

	var durable, del, exclusive, noWait bool

	if _, ok := attributes["durable"]; ok {
		durable = attributes["durable"].(bool)
		delete(attributes, "durable")
	} else {
		durable = true
	}

	if _, ok := attributes["delete"]; ok {
		del = attributes["delete"].(bool)
		delete(attributes, "delete")
	} else {
		del = false
	}

	if _, ok := attributes["exclusive"]; ok {
		exclusive = attributes["exclusive"].(bool)
		delete(attributes, "exclusive")
	} else {
		exclusive = false
	}

	if _, ok := attributes["no-wait"]; ok {
		noWait = attributes["no-wait"].(bool)
		delete(attributes, "no-wait")
	} else {
		noWait = false
	}

	fmt.Printf("durable: %v, delete: %v, name: %s, exclusive: %v, no-wait: %v, attributes: %v\n", durable, del, name, exclusive, noWait, attributes)
	_, err := _amqp.Channel.QueueDeclare(
		name,       // name
		durable,    // durable
		del,        // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		attributes, // arguments
	)

	return err

}

// ReceiveMessages is an AMQP convenience method to receive messages from a given queue
func (_amqp *AMQP) ReceiveMessages(queue string) (<-chan Message, error) {

	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to declare queue %s: %s", queue, err.Error())
	// }

	output, err := _amqp.Channel.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return nil, err
	}

	ret := make(chan Message)
	go func() {
		for message := range output {
			ret <- Message{message.MessageId, message.CorrelationId, message.Body}
		}
	}()

	return ret, nil
}
