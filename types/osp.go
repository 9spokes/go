package types

//OSPDefinition is an Online Service Provider definition
type OSPDefinition struct {
	Name                 string            `json:"name"`
	Credentials          map[string]string `json:"credentials"`
	NotificationSettings []struct {
		Metric string   `json:"metric"`
		Period []string `json:"period"`
	} `json:"notification_settings"`
}
