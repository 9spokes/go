package types

//OSP is an Online Service Provider definition
type OSP struct {
	Name        string            `json:"osp"`
	NewETL      bool              `json:"new_etl"`
	Credentials map[string]string `json:"credentials,omitempty"`
	Tiles       []OSPTile         `json:"tiles"`
	Proxy       string            `json:"proxy,omitempty"`
	Unique      bool              `json:"unique,omitempty"`
	Usage       []string          `json:"usage,omitempty"`
}

// OSPTile is a tile definition for an OSP
type OSPTile struct {
	Name                 string `json:"name"`
	NotificationSettings []struct {
		Metric string   `json:"metric,omitempty"`
		Period []string `json:"period,omitempty"`
	} `json:"notification_settings,omitempty"`
}
