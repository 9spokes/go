package companies

import (
	"time"

	"github.com/9spokes/go/types"
)

//Company is a struct that defines a 9 Spokes company entity
type Company struct {
	ID                string        `json:"id" bson:"company"`
	Name              string        `json:"name" bson:"name"`
	Industry          string        `json:"industry,omitempty" bson:"industry,omitempty"`
	Industries        []Record      `json:"industries"`
	Location          types.Place   `json:"location"`
	Type              Record        `json:"type"`
	Status            Record        `json:"status"`
	Entity            string        `json:"entity,omitempty" bson:"entity,omitempty"`
	Phone             string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Users             []Users       `json:"users"`
	TimeZoneOffset    int64         `json:"timeZoneOffset,omitempty" bson:"timeZoneOffset,omitempty"`
	BusinessHours     BusinessHours `json:"businessHours"`
	WorkingHoursStart string        `json:"workingHoursStart,omitempty"`
	WorkingHoursEnd   string        `json:"workingHoursEnd,omitempty"`
	StartTime         string        `json:"startTime,omitempty"`
	EndTime           string        `json:"endTime,omitempty"`
	Created           time.Time     `json:"created"`
	Updated           time.Time     `json:"updated"`
	Details           *Details      `json:"details,omitempty" bson:"details,omitempty"`
	Extras            Extras        `json:"extras,omitempty" bson:"extras,omitempty"`
	Size              string        `json:"size,omitempty" bson:"size,omitempty"`
}

// Users represent the user details of the company
type Users struct {
	User     string `json:"user"`
	Role     string `json:"role"`
	Position string `json:"position"`
}

// BusinessHours represent the working days & hours for a company
type BusinessHours struct {
	DaysFrom  int `json:"daysFrom"`
	DaysTo    int `json:"daysTo"`
	HoursFrom int `json:"hoursFrom"`
	HoursTo   int `json:"hoursTo"`
}

// Address represents a company address
type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	AddressLine4 string `json:"addressLine4"`
	PostCode     string `json:"postCode"`
	Country      string `json:"country"`
	CareOf       string `json:"careOf"`
	Type         string `json:"type"`
}

// AnnualReturn represents a company's annual return details
type AnnualReturn struct {
	FilingMonth int       `json:"filingMonth"`
	LastFiled   time.Time `json:"lastFiled"`
}

// Director represents a company director
type Director struct {
	Name        string    `json:"name"`
	AppointedAt time.Time `json:"appointedAt"`
	Status      string    `json:"status"`
}

// ShareAllocation represents a company share allocation
type ShareAllocation struct {
	Type       string `json:"type"`
	Allocation int    `json:"allocation"`
	Name       string `json:"shareholder"`
}

// Record represents a code and description touple
type Record struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Summary is a struct that defines an external Company Registry company short summary
type Summary struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     Record   `json:"status"`
	Type       Record   `json:"type"`
	Industries []Record `json:"industries"`
}

type Extras struct {
	AnnualTurnover float64 `json:"annualTurnover,omitempty" bson:"annualTurnover,omitempty"`
	Size           int64   `json:"size,omitempty" bson:"size,omitempty"`
}

// Details is a struct that defines a company's details as retrieved from an external Company Registry
type Details struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Status            string                 `json:"status"`
	Type              string                 `json:"type"`
	IncorporationDate string                 `json:"incorporationDate"`
	Emails            []string               `json:"emails"`
	Addresses         []Address              `json:"addresses"`
	Industry          []Record               `json:"industry"`
	AnnualReturn      AnnualReturn           `json:"annualReturn"`
	TotalShares       int64                  `json:"totalShares"`
	ShareAllocation   []ShareAllocation      `json:"shareAllocation"`
	Directors         []Director             `json:"directors"`
	Extras            map[string]interface{} `json:"extras"`
	Modified          time.Time              `json:"modified"`
}
