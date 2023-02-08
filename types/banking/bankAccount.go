package banking

type BankAccount struct {
	Id               string          `json:"accountId" bson:"account_id"`
	Type             BankAccountType `json:"type" bson:"type"`
	Number           string          `json:"accountNumber" bson:"account_number"`
	NumberMasked     string          `json:"accountNumberMasked" bson:"account_number_masked"`
	DisplayName      string          `json:"displayName" bson:"display_name"`
	Currency         string          `json:"currency" bson:"currency"`
	CurrentBalance   float64         `json:"currentBalance" bson:"currentBalance"`
	AvailableBalance float64         `json:"availableBalance" bson:"availableBalance"`
}
