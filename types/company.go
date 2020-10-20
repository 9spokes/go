package types

import "time"

//Company is a struct that defines a 9 Spokes company entity
type Company struct {
	ID                string                 `json:"id" bson:"company"`
	Name              string                 `json:"name" bson:"name"`
	Industry          string                 `json:"industry,omitempty"`
	Industries        []Record               `json:"industries"`
	Location          CompanyLocation        `json:"location"`
	Type              Record                 `json:"type"`
	Status            Record                 `json:"status"`
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
	AddressLine1 string `json:"address-line1"`
	AddressLine2 string `json:"address-line2"`
	AddressLine3 string `json:"address-line3"`
	AddressLine4 string `json:"address-line4"`
	PostCode     string `json:"post-code"`
	Country      string `json:"country"`
	CareOf       string `json:"care-of"`
	Type         string `json:"type"`
}

// CompanyAnnualReturn represents a company's annual return details
type CompanyAnnualReturn struct {
	FilingMonth int       `json:"filing-month"`
	LastFiled   time.Time `json:"last-filed"`
}

// CompanyDirector represents a company director
type CompanyDirector struct {
	Name        string    `json:"name"`
	AppointedAt time.Time `json:"appointed-at"`
	Status      string    `json:"status"`
}

// CompanyShareAllocation represents a company share allocation
type CompanyShareAllocation struct {
	Type       string `json:"type"`
	Allocation int    `json:"allocation"`
	Name       string `json:"shareholder"`
}

// Record represents a code and description touple
type Record struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// CompanySummary is a struct that defines an external Company Registry company short summary
type CompanySummary struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     Record   `json:"status"`
	Type       Record   `json:"type"`
	Industries []Record `json:"industries"`
}

// CompanyDetails is a struct that defines a company's details as retrieved from an external Company Registry
type CompanyDetails struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	Status            string                   `json:"status"`
	Type              string                   `json:"type"`
	IncorporationDate string                   `json:"incorporation-date"`
	Emails            []string                 `json:"emails"`
	Addresses         []CompanyAddress         `json:"addresses"`
	Industry          []Record                 `json:"industry"`
	AnnualReturn      CompanyAnnualReturn      `json:"annual-return"`
	TotalShares       int64                    `json:"total-shares"`
	ShareAllocation   []CompanyShareAllocation `json:"share-allocation"`
	Directors         []CompanyDirector        `json:"directors"`
	Extras            map[string]interface{}   `json:"extras"`
	Modified          time.Time                `json:"modified"`
}
