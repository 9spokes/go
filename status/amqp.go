package status

import (
	"github.com/streadway/amqp"
)

//ValidateAMQP validates an AMQP connection
func ValidateAMQP(amqp *amqp.Connection) bool {
	return !amqp.IsClosed()
}
