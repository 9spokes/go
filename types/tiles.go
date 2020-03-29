package types

// ListTile is a 9 Spokes V2 "List" tile data format
type ListTile struct {
	SubHeader struct {
		ValueLeft float64 `json:"valueLeft,omitempty"`
	} `json:"subheader,omitempty"`
	List []ListTileEntry `json:"list,omitempty"`
}

// ListTileEntry is an entry in a ListTile
type ListTileEntry struct {
	Label     string  `json:"label,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Indicator string  `json:"indicator,omitempty"`
	Left      string  `json:"left,omitempty"`
	Right     string  `json:"right,omitempty"`
	IsSubRow  bool    `json:"isSubRow,omitempty"`
	Footer    struct {
		Label string  `json:"label,omitempty"`
		Value float64 `json:"value,omitempty"`
	} `json:"footer,omitempty"`
}

// GraphTile is a 9 Spokes V2 area chart tile data format
type GraphTile struct {
	XUnit  string        `json:"xUnit,omitempty"`
	YUnit  float64       `json:"yUnit,omitempty"`
	Labels []interface{} `json:"labels,omitempty"`
	Series []struct {
		Key  string    `json:"key,omitempty"`
		Data []float64 `json:"data,omitempty"`
	} `json:"series,omitempty"`
}
