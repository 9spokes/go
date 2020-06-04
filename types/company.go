package types

//Company is a struct that defines a 9 Spokes company entity
type Company struct {
	ID                string                 `json:"id" bson:"connection"`
	Name              string                 `json:"name" bson:"name"`
	Industry          string                 `json:"industry"`
	Phone             string                 `json:"phone"`
	Users             []string               `json:"users"`
	TimeZoneOffset    int64                  `json:"timeZoneOffset"`
	WorkingHoursStart string                 `json:"workingHoursStart"`
	WorkingHoursEnd   string                 `json:"workingHoursEnd"`
	Extras            map[string]interface{} `json:"extras"`
}
