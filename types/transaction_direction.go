package types

//go:generate go run github.com/dmarkham/enumer@latest -type=TransactionDirection -json
type TransactionDirection int

const (
	Inbound TransactionDirection = iota + 1
	Outbound
)
