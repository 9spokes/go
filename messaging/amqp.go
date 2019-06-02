package messaging

import (
	"fmt"
	"strconv"

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

	var exchange string
	var mandatory, immediate bool
	var priority uint8
	var expiration string

	if message.Options == nil {
		message.Options = make(map[string]interface{})
	}

	if _, ok := message.Options["exchange"]; ok {
		exchange = message.Options["exchange"].(string)
		delete(message.Options, "exchange")
	} else {
		exchange = ""
	}

	if _, ok := message.Options["mandatory"]; ok {
		mandatory = message.Options["mandatory"].(bool)
		delete(message.Options, "mandatory")
	} else {
		mandatory = false
	}

	if _, ok := message.Options["immediate"]; ok {
		immediate = message.Options["immediate"].(bool)
		delete(message.Options, "immediate")
	} else {
		immediate = false
	}

	if _, ok := message.Options["priority"]; ok {
		priority = message.Options["priority"].(uint8)
		delete(message.Options, "priority")
	} else {
		priority = 0
	}

	if _, ok := message.Options["x-message-ttl"]; ok {
		expiration = strconv.FormatInt(message.Options["x-message-ttl"].(int64), 10)
		delete(message.Options, "x-message-ttl")
	} else {
		expiration = ""
	}

	err := _amqp.Channel.Publish(
		exchange,  // exchange
		queue,     // routing key
		mandatory, // mandatory
		immediate, // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          message.Body,
			CorrelationId: message.CorrelationID,
			Headers:       message.Options,
			Priority:      priority,
			Expiration:    expiration,
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
func (_amqp *AMQP) ReceiveMessages(queue string, opt map[string]interface{}) (<-chan Message, error) {

	var consumer string
	var autoAck, exclusive, noLocal, noWait bool

	if _, ok := opt["consumer"]; ok {
		consumer = opt["consumer"].(string)
		delete(opt, "consumer")
	} else {
		consumer = ""
	}

	if _, ok := opt["auto-ack"]; ok {
		autoAck = opt["auto-ack"].(bool)
		delete(opt, "auto-ack")
	} else {
		autoAck = false
	}

	if _, ok := opt["exclusive"]; ok {
		exclusive = opt["exclusive"].(bool)
		delete(opt, "exclusive")
	} else {
		exclusive = false
	}

	if _, ok := opt["no-local"]; ok {
		noLocal = opt["no-local"].(bool)
		delete(opt, "no-local")
	} else {
		noLocal = true
	}

	if _, ok := opt["no-wait"]; ok {
		noWait = opt["no-wait"].(bool)
		delete(opt, "no-wait")
	} else {
		noWait = false
	}

	output, err := _amqp.Channel.Consume(
		queue,     // queue
		consumer,  // consumer
		autoAck,   // auto-ack
		exclusive, // exclusive
		noLocal,   // no-local
		noWait,    // no-wait
		opt,       // args
	)

	if err != nil {
		return nil, err
	}

	ret := make(chan Message)
	go func() {
		opt := make(map[string]interface{})

		for message := range output {
			opt["timestamp"] = message.Timestamp
			opt["priority"] = message.Priority
			opt["messageCount"] = message.MessageCount
			opt["exchange"] = message.Exchange
			opt["routingKey"] = message.RoutingKey
			opt["redelivered"] = message.Redelivered

			ret <- Message{ID: message.MessageId, CorrelationID: message.CorrelationId, Body: message.Body, Options: opt}
		}
	}()

	return ret, nil
}
