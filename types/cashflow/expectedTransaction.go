package cashflow

import (
	"time"

	"github.com/9spokes/go/types"
)

type ExpectedTransaction struct {
	Id              string                     `json:"transactionId" bson:"transaction_id"`
	User            string                     `json:"userId" bson:"user_id"`
	TransactionDate time.Time                  `json:"transactionDate" bson:"transaction_date"`
	Counterparty    string                     `json:"counterparty" bson:"counterparty"`
	Amount          float64                    `json:"amount" bson:"amount"`
	Currency        string                     `json:"currency" bson:"currency"`
	Direction       types.TransactionDirection `json:"transactionDirection" bson:"transaction_direction"`
	AccountID       string                     `json:"accountId" bson:"account_id"`
}
