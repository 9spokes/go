package types

// Holds an OSP's configuration
type OSP struct {
	Name            string            `json:"osp"`
	NewETL          bool              `json:"new_etl"`
	ExtractOnSignin bool              `json:"extract_on_signin"`
	Credentials     map[string]string `json:"credentials,omitempty"`
	Tiles           []OSPTile         `json:"tiles"`
	Providers       []string          `json:"providers,omitempty"`
	Proxy           string            `json:"proxy,omitempty"`
	Unique          bool              `json:"unique,omitempty"`
	Usage           []string          `json:"usage,omitempty"`
	TLS             map[string]string `json:"tls,omitempty"`
}

// OSPTile is a tile definition for an OSP
type OSPTile struct {
	Name                 string `json:"name"`
	NotificationSettings []struct {
		Metric string   `json:"metric,omitempty"`
		Period []string `json:"period,omitempty"`
	} `json:"notification_settings,omitempty"`
}
