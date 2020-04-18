package types

import "time"

// Document is a generic map
type Document map[string]interface{}

// Connection represents a Connection document object as stored in the database
type Connection struct {
	ID            string    `json:"id" bson:"connection"`
	Credentials   Document  `bson:"credentials" json:"credentials"`
	Demo          bool      `bson:"demo" json:"demo" default:"false"`
	Token         Document  `bson:"token" json:"token"`
	Settings      Document  `bson:"settings" json:"settings"`
	User          string    `bson:"user" json:"user"`
	Configuration Document  `bson:"config" json:"config"`
	OSP           string    `bson:"osp" json:"osp"`
	Created       time.Time `bson:"created" json:"created"`
	Modified      time.Time `bson:"modified" json:"modified"`
	Status        string    `bson:"status" default:"NEW" json:"status"`
}

// ConnectionSummary is a short-form connection object as returned by the token service, it excludes sensitive info and is meant as a summary
type ConnectionSummary struct {
	ID           string      `json:"id,omitempty"`
	Created      time.Time   `json:"created" bson:"created,omitempty"`
	Modified     time.Time   `json:"modified" bson:"modified,omitempty"`
	OSP          string      `json:"osp" bson:"osp,omitempty"`
	Status       string      `json:"status" bson:"status,omitempty"`
	AuthorizeURL string      `json:"authorize_url,omitempty"`
	Settings     interface{} `json:"settings"`
}
