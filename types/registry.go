package types

import "time"

// CompanyAddress represents a company address
type CompanyAddress struct {
	AdressLine1 string `json:"addressLine1"`
	AdressLine2 string `json:"addressLine2"`
	AdressLine3 string `json:"addressLine3"`
	AdressLine4 string `json:"addressLine4"`
	PostCode    string `json:"postCode"`
	Country     string `json:"country"`
	CareOf      string `json:"careOf"`
	Type        string `json:"type"`
}

// CompanyAnnualReturn represents a company's annual return details
type CompanyAnnualReturn struct {
	FilingMonth int       `json:"filingMonth"`
	LastFiled   time.Time `json:"lastFiled"`
}

// CompanyDirector represents a company director
type CompanyDirector struct {
	Name            string    `json:"name"`
	AppointmentDate time.Time `json:"appointmentDate"`
	Status          string    `json:"status"`
}

// ShareAllocation represents a company share allocation
type ShareAllocation struct {
	Type       string `json:"type"`
	Allocation int    `json:"allocation"`
	Name       string `json:"shareholderName"`
}

// IndustryClassification represents the company's industry classification
type IndustryClassification struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// CompanySummary is a struct that defines an external Company Registry company short summary
type CompanySummary struct {
	ID         string `json:"id"`
	Entityname string `json:"name"`
	Status     string `json:"status"`
	Type       string `json:"type"`
}

// CompanyDetails is a struct that defines a company's details as retrieved from an external Company Registry
type CompanyDetails struct {
	ID                      string                   `json:"id"`
	Entityname              string                   `json:"name"`
	Status                  string                   `json:"status"`
	Type                    string                   `json:"type"`
	IncorporationDate       string                   `json:"incorporationDate"`
	Emails                  []string                 `json:"emails"`
	Adressess               []CompanyAddress         `json:"adresses"`
	IndustryClassifications []IndustryClassification `json:"industryClassifications"`
	AnnualReturn            CompanyAnnualReturn      `json:"annualReturn"`
	TotalShares             int                      `json:"totalShares"`
	ShareAllocation         []ShareAllocation        `json:"shareAllocation"`
	Directors               []CompanyDirector        `json:"directors"`
	Extras                  map[string]interface{}   `json:"extras"`
}
