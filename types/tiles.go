package types

// ListTileSubheader represents a list tile subheader
type ListTileSubheader struct {
	Title 	   string `json:"title,omitempty"`
	LabelLeft  string `json:"labelLeft,omitempty"`
	ValueLeft  string `json:"valueLeft,omitempty"`
	LabelRight string `json:"labelRight,omitempty"`
	ValueRight string `json:"valueRight,omitempty"`
}

// ListTileSubheaderDropDownOption represents an entry in the list tile
// subheader dropdown
type ListTileSubheaderDropDownOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ListTileSubheaderDropDownOpt struct {
	Title 	   string `json:"title,omitempty"`
	LabelLeft  string `json:"labelLeft,omitempty"`
	ValueLeft  []ListTileSubheaderDropDownOption `json:"valueLeft,omitempty"`
	LabelRight string `json:"labelRight,omitempty"`
	ValueRight string `json:"valueRight,omitempty"`
}

// ListTile is a 9 Spokes V2 "List" tile data format
type ListTile struct {
	SubHeader ListTileSubheader `json:"subheader,omitempty"`
	List      []ListTileEntry   `json:"list,omitempty"`
	SyncedAt  string            `json:"lastSyncAt,omitempty"`
}

// The same as the ListTile but the subheader has Title property
type ListTileWithSubheaderTitle struct {
	SubHeader ListTileSubheaderDropDownOpt `json:"subheader,omitempty"`
	List      []ListTileEntry   `json:"list,omitempty"`
	SyncedAt  string            `json:"lastSyncAt,omitempty"`
}

// ListTileEntry is an entry in a ListTile
type ListTileEntry struct {
	Label     string `json:"label,omitempty"`
	Value     string `json:"value,omitempty"`
	Indicator string `json:"indicator,omitempty"`
	Left      string `json:"left,omitempty"`
	Right     string `json:"right,omitempty"`
	IsSubRow  bool   `json:"isSubRow,omitempty"`
	Footer    struct {
		Label string `json:"label,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"footer,omitempty"`
}

// GraphTile is a 9 Spokes V2 area chart tile data format
type GraphTile struct {
	XUnit  string   `json:"xUnit,omitempty"`
	YUnit  string   `json:"yUnit,omitempty"`
	Labels []string `json:"labels,omitempty"`
	Series []struct {
		Key  string    `json:"key,omitempty"`
		Data []float64 `json:"data,omitempty"`
	} `json:"series,omitempty"`
	SyncedAt string `json:"lastSyncAt,omitempty"`
}

// CompositeListTile is a 9 Spokes V2 special "List" tile data format
// It's made up of multiple ListTile blocks so the tile can switch
// between them
type CompositeListTile struct {
	GroupedData []GroupedDataEntry `json:"groupedData,omitempty"`
}

// GroupedDataEntry represenst a data set in a composite list tile
type GroupedDataEntry struct {
	Key  string   `json:"key,omitempty"`
	Data ListTile `json:"data,omitempty"`
}
