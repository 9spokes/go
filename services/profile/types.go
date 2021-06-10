package profile

import "time"

// Profile is a struct representing the user's Tracker profile
type Profile struct {
	User            string                 `bson:"user" json:"user" scope:"private"`
	FirstName       string                 `bson:"firstname" json:"firstname,omitempty" scope:"public,private"`
	LastName        string                 `bson:"lastname" json:"lastname,omitempty" scope:"public,private"`
	Email           string                 `bson:"email" json:"email,omitempty" scope:"private"`
	Created         time.Time              `bson:"created" json:"-" scope:"private"`
	LastLogin       time.Time              `bson:"lastLogin" json:"lastLogin" scope:"private"`
	EmailVerified   bool                   `bson:"emailVerified" json:"emailVerified,omitempty" scope:"private"`
	Company         string                 `bson:"company" json:"-" scope:"private"`
	Terms           string                 `bson:"terms" json:"terms,omitempty" scope:"private"`
	Demo            bool                   `bson:"demo" json:"demo,omitempty" scope:"private"`
	StorefrontToken string                 `bson:"storefrontToken"  json:"storefrontToken,omitempty"`
	LoginCount      uint64                 `bson:"loginCount" json:"loginCount" scope:"private"`
	Extras          map[string]interface{} `bson:"extras,omitempty" json:"extras,omitempty" scope:"private"`
}
