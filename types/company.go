package types

import "time"

//Company is a struct that defines a 9 Spokes company entity
type Company struct {
	ID                string                 `json:"id" bson:"company"`
	Name              string                 `json:"name" bson:"name"`
	Industry          string                 `json:"industry,omitempty"`
	Industries        []CompanyRecord        `json:"industries"`
	Location          CompanyLocation        `json:"location"`
	Type              CompanyRecord          `json:"type"`
	Status            CompanyRecord          `json:"status"`
	Entity            string                 `json:"entity,omitempty"`
	Phone             string                 `json:"phone,omitempty"`
	Users             []string               `json:"users"`
	TimeZoneOffset    int64                  `json:"timeZoneOffset,omitempty"`
	WorkingHoursStart string                 `json:"workingHoursStart,omitempty"`
	WorkingHoursEnd   string                 `json:"workingHoursEnd,omitempty"`
	Created           time.Time              `json:"created"`
	Updated           time.Time              `json:"updated"`
	Details           *CompanyDetails        `json:"details,omitempty"`
	Extras            map[string]interface{} `json:"extras,omitempty"`
}

// CompanyLocation contains the company's locality details
type CompanyLocation struct {
	Country  string `json:"country"`
	Timezone int    `json:"timezone"`
}

// CompanyAddress represents a company address
type CompanyAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	AddressLine4 string `json:"addressLine4"`
	PostCode     string `json:"postCode"`
	Country      string `json:"country"`
	CareOf       string `json:"careOf"`
	Type         string `json:"type"`
}

// CompanyAnnualReturn represents a company's annual return details
type CompanyAnnualReturn struct {
	FilingMonth int       `json:"filingMonth"`
	LastFiled   time.Time `json:"lastFiled"`
}

// CompanyDirector represents a company director
type CompanyDirector struct {
	Name        string    `json:"name"`
	AppointedAt time.Time `json:"appointedAt"`
	Status      string    `json:"status"`
}

// CompanyShareAllocation represents a company share allocation
type CompanyShareAllocation struct {
	Type       string `json:"type"`
	Allocation int    `json:"allocation"`
	Name       string `json:"shareholder"`
}

// CompanyRecord represents a code and description touple
type CompanyRecord struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// CompanySummary is a struct that defines an external Company Registry company short summary
type CompanySummary struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Status     CompanyRecord   `json:"status"`
	Type       CompanyRecord   `json:"type"`
	Industries []CompanyRecord `json:"industries"`
}

// CompanyDetails is a struct that defines a company's details as retrieved from an external Company Registry
type CompanyDetails struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	Status            string                   `json:"status"`
	Type              string                   `json:"type"`
	IncorporationDate string                   `json:"incorporationDate"`
	Emails            []string                 `json:"emails"`
	Addresses         []CompanyAddress         `json:"addresses"`
	Industry          []CompanyRecord          `json:"industry"`
	AnnualReturn      CompanyAnnualReturn      `json:"annualReturn"`
	TotalShares       int64                    `json:"totalShares"`
	ShareAllocation   []CompanyShareAllocation `json:"shareAllocation"`
	Directors         []CompanyDirector        `json:"directors"`
	Extras            map[string]interface{}   `json:"extras"`
	Modified          time.Time                `json:"modified"`
}
