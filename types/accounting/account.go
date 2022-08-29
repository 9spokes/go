package accounting

type accountCategory string

const (
	Asset     accountCategory = "assets"
	Liability accountCategory = "liability"
	Expense   accountCategory = "expense"
	Revenue   accountCategory = "revenue"
	Equity    accountCategory = "equity"
)

type accountType string

const (
	AccountsPayable    accountType = "current_accounts_payable"
	AccountsReceivable accountType = "current_accounts_receivable"
	Bank               accountType = "bank"
	Current            accountType = "current"
	Depreciation       accountType = "depreciation"
	Fixed              accountType = "fixed"
	Inventory          accountType = "inventory"
	Other              accountType = "other"
	Overheads          accountType = "overheads"
	Owners             accountType = "owners"
	Payroll            accountType = "payroll"
	Prepayments        accountType = "prepayments"
	RetainedEarnings   accountType = "retained-earnings"
	Sales              accountType = "sales"
	Tax                accountType = "tax"
	Term               accountType = "term"
)

type accountStatus string

const (
	Active   accountStatus = "ACTIVE"
	Inactive accountStatus = "INACTIVE"
)

type Account struct {
	ID       string          `json:"account_id,omitempty" bson:"account_id,omitempty"`
	Active   bool            `json:"account_active,omitempty" bson:"account_active,omitempty"` // TODO: remove Status or Active
	Balance  float64         `json:"account_balance" bson:"account_balance"`
	Currency string          `json:"account_currency,omitempty" bson:"account_currency,omitempty"`
	Name     string          `json:"account_name,omitempty" bson:"account_name,omitempty"`
	Type     accountType     `json:"account_type,omitempty" bson:"account_type,omitempty"`
	Category accountCategory `json:"account_category,omitempty" bson:"account_category,omitempty"`
	Status   accountStatus   `json:"account_status,omitempty" bson:"account_status,omitempty"` // TODO: remove Status or Active
	Number   string          `json:"account_num,omitempty" bson:"account_num,omitempty"`
}
