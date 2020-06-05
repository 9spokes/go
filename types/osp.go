package types

//OSP is an Online Service Provider definition
type OSP struct {
	Name        string            `json:"osp"`
	Credentials map[string]string `json:"credentials,omit_empty"`
	Tiles       []OSPTile         `json:"tiles"`
}

// OSPTile is a tile definition for an OSP
type OSPTile struct {
	Name                 string `json:"name'`
	NotificationSettings []struct {
		Metric string   `json:"metric,omit_empty"`
		Period []string `json:"period,omit_empty"`
	} `json:"notification_settings,omit_empty"`
}
