package banking

//go:generate go run github.com/dmarkham/enumer@latest -type=BankAccountType -json
type BankAccountType int

const (
	Business BankAccountType = iota + 1
	Personal
)
