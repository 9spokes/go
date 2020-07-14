package types

import "time"

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

// CompanyIndustry represents the company's industry classification
type CompanyIndustry struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// CompanySummary is a struct that defines an external Company Registry company short summary
type CompanySummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
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
	Industry          []CompanyIndustry        `json:"industry"`
	AnnualReturn      CompanyAnnualReturn      `json:"annual-return"`
	TotalShares       int64                    `json:"total-shares"`
	ShareAllocation   []CompanyShareAllocation `json:"share-allocation"`
	Directors         []CompanyDirector        `json:"directors"`
	Extras            map[string]interface{}   `json:"extras"`
}
