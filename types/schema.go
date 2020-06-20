package types

// MarketingAnalyticsDateRange represents a date range used to tag a 9 Spokes data record
type MarketingAnalyticsDateRange struct {
	Start string `json:"event_start"`
	End   string `json:"event_end"`
}

// MarketingAnalytics is the standard 9 Spokes Marketing Analytics data format
type MarketingAnalytics struct {
	User       string                        `json:"user"`
	Connection string                        `json:"connection"`
	Company    string                        `json:"company"`
	OSP        string                        `json:"osp"`
	Type       string                        `json:"type"`
	Updated    string                        `json:"updated"`
	Period     string                        `json:"period"`
	Index      string                        `json:"index"`
	Cycle      string                        `json:"cycle"`
	Datasource string                        `json:"datasource"`
	DateRanges []MarketingAnalyticsDateRange `json:"date_ranges"`
	Data       interface{}                   `json:"data"`
}
