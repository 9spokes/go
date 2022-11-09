package banking

import (
	"time"

	"github.com/9spokes/go/types"
)

type BankTransaction struct {
	Id                   string                `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	Date                 time.Time             `json:"transaction_date,omitempty" bson:"transaction_date,omitempty"`
	Description          string                `json:"description,omitempty" bson:"description,omitempty"`
	Reconciled           bool                  `json:"reconciled,omitempty" bson:"reconciled,omitempty"`
	Amount               float64               `json:"amount,omitempty" bson:"amount,omitempty"`
	Balance              float64               `json:"balance,omitempty" bson:"balance,omitempty"`
	Currency             string                `json:"currency,omitempty" bson:"currency,omitempty"`
	AccountID            string                `json:"account_id,omitempty" bson:"account_id,omitempty"`
	AccountNumber        string                `json:"account_num,omitempty" bson:"account_num,omitempty"`
	AccountNumberDisplay string                `json:"account_num_display,omitempty" bson:"account_num_display,omitempty"`
	AccountName          string                `json:"account_name,omitempty" bson:"account_name,omitempty"`
	LedgerID             string                `json:"ledger_id,omitempty" bson:"ledger_id,omitempty"`
	Status               string                `json:"status,omitempty" bson:"status,omitempty"`
	ValueType            string                `json:"value_type,omitempty" bson:"value_type,omitempty"`
	Payee                string                `json:"payee,omitempty" bson:"payee,omitempty"`
	Type                 types.TransactionType `json:"transaction_type,omitempty" bson:"transaction_type,omitempty"`
}
