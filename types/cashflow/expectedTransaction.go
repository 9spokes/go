package types

import (
	"time"
)

type ExpectedTransaction struct {
	Id              string               `json:"transactionId" bson:"transaction_id"`
	User            string               `json:"userId" bson:"user_id"`
	TransactionDate time.Time            `json:"transactionDate" bson:"transaction_date"`
	Counterparty    string               `json:"counterparty" bson:"counterparty"`
	Type            TransactionType      `json:"transactionType" bson:"transaction_type"`
	Amount          float64              `json:"amount" bson:"amount"`
	Currency        string               `json:"currency" bson:"currency"`
	Direction       TransactionDirection `json:"transactionDirection" bson:"transaction_direction"`
	AccountID       string               `json:"accountId" bson:"account_id"`
}

//go:generate go run github.com/dmarkham/enumer -type=TransactionType -json
type TransactionType int

const (
	Uncategorised TransactionType = iota
	Cash
	Cheque
	Purchase
	BillPayment
	AutomaticPayment
	Fee
	Adjustment
)

//go:generate go run github.com/dmarkham/enumer -type=TransactionDirection -json
type TransactionDirection int

const (
	Inbound TransactionDirection = iota
	Outbound
)
