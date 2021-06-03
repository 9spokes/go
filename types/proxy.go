package types

//CodatCompany struct contains information of a 9Spokes user registering to Codat
type CodatCompany struct {
	User      string `json:"user,omitempty" bson:"user,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	CompanyID string `json:"company_id,omitempty" bson:"company_id,omitempty"`
	Alert     bool   `json:"alert" bson:"alert"`
}
