package types

import "time"

// Document is a generic map
type Document map[string]interface{}

// Connection represents a Connection document object as stored in the database
type Connection struct {
	ID            string    `json:"id" bson:"connection"`               // The unique connection ID
	Platform      string    `json:"platform" bson:"platform"`           // Either "tracker" or "odp"
	Credentials   Document  `bson:"credentials" json:"credentials"`     // A map of key/value pairs representing OSP credentials
	Demo          bool      `bson:"demo" json:"demo" default:"false"`   // Whether this is a demo connection or not, determines how to render tiles
	Token         Document  `bson:"token" json:"token"`                 // Contains the encrypted and decrypted access & refresh tokens
	Settings      Document  `bson:"settings" json:"settings"`           // Contains app-specific settings for this connection
	User          string    `bson:"user" json:"user"`                   // The UUID of the user that created this connection
	Configuration Document  `bson:"config" json:"config"`               // Contains additioanl OSP-specific configuration
	OSP           string    `bson:"osp" json:"osp"`                     // The App
	Proxy         string    `bson:"proxy" json:"proxy,omitempty"`       // Whether this connection shoudl be proxied through a 3rd party of direct
	Usage         []string  `bson:"usage" json:"usage,omitempty"`       //
	Company       string    `bson:"company" json:"company"`             // The UUID of the company that owns this connection
	Created       time.Time `bson:"created" json:"created"`             // The RFC3339 creation date of this connection
	Modified      time.Time `bson:"modified" json:"modified"`           // The RFC3339 modification date of this connection
	Status        string    `bson:"status" default:"NEW" json:"status"` // Either ACTIVE, NOT_CONNECTED, or NEW
}

// ConnectionSummary is a short-form connection object as returned by the token service, it excludes sensitive info and is meant as a summary
type ConnectionSummary struct {
	ID           string    `json:"id,omitempty"`
	Created      time.Time `json:"created" bson:"created,omitempty"`
	Modified     time.Time `json:"modified" bson:"modified,omitempty"`
	OSP          string    `json:"osp" bson:"osp,omitempty"`
	Status       string    `json:"status" bson:"status,omitempty"`
	AuthorizeURL string    `json:"authorize_url,omitempty"`
	Settings     Document  `json:"settings"`
}
