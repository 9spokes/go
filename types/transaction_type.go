package types

//go:generate go run github.com/dmarkham/enumer@latest -type=TransactionType -json
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
	Transfer
	Interest
	Expected
	Projected TransactionType = Expected
)
