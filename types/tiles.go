package types

// ListTileSubheader represents a list tile subheader
type ListTileSubheader struct {
	Title      string `json:"title,omitempty"`
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
	Title       string                            `json:"title,omitempty"`
	LabelLeft   string                            `json:"labelLeft,omitempty"`
	ValueLeft   []ListTileSubheaderDropDownOption `json:"valueLeft,omitempty"`
	LabelRight  string                            `json:"labelRight,omitempty"`
	ValueRight  string                            `json:"valueRight,omitempty"`
	MultiSelect bool                              `json:"multiSelect,omitempty"`
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
	List      []ListTileEntry              `json:"list,omitempty"`
	Graph     GraphTile                    `json:"graph,omitempty"`
	SyncedAt  string                       `json:"lastSyncAt,omitempty"`
}

type ListTileFooter struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

// ListTileEntry is an entry in a ListTile
type ListTileEntry struct {
	Label     string         `json:"label,omitempty"`
	Value     string         `json:"value,omitempty"`
	Indicator string         `json:"indicator,omitempty"`
	Left      string         `json:"left,omitempty"`
	Right     string         `json:"right,omitempty"`
	IsSubRow  bool           `json:"isSubRow,omitempty"`
	Direction string         `json:"direction,omitempty"`
	Footer    ListTileFooter `json:"footer,omitempty"`
	Icon      string         `json:"icon,omitempty"`
	Event     *ListTileEvent `json:"event,omitempty"`
}

type ListTileEvent struct {
	OnClick *ListTileClickEvent `json:"onClick,omitempty"`
}

type ListTileClickEvent struct {
	Action string `json:"action,omitempty"`
	URL    string `json:"url,omitempty"`
	Value  string `json:"value,omitempty"`
}

// GraphTile is a 9 Spokes V2 area chart tile data format
type GraphTile struct {
	XUnit     string             `json:"xUnit,omitempty"`
	YUnit     string             `json:"yUnit,omitempty"`
	LabelData GraphTileLabelData `json:"labelData,omitempty"`
	Series    []GraphTileSeries  `json:"series,omitempty"`
	XGroups   []string           `json:"xGroups,omitempty"`
	Values    []GraphTileData    `json:"values,omitempty"`
	SyncedAt  string             `json:"lastSyncAt,omitempty"`
}

type GraphTileLabelData struct {
	Labels    []string `json:"labels,omitempty"`
	FormatKey string   `json:"formatKey,omitempty"`
}

type GraphTileSeries struct {
	Key  string    `json:"key,omitempty"`
	Data []float64 `json:"data,omitempty"`
}
type GraphTileData struct {
	Key      string           `json:"key,omitempty"`
	Category string           `json:"category,omitempty"`
	Values   []GraphTileValue `json:"values,omitempty"`
}

type GraphTileValue struct {
	Group string  `json:"group,omitempty"`
	Stack string  `json:"stack,omitempty"`
	Value float64 `json:"value"`
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

type CompositeListTileWithSubTitle struct {
	GroupedData []GroupedDataEntryWithSubTitle `json:"groupedData,omitempty"`
}

type GroupedDataEntryWithSubTitle struct {
	Key  string                     `json:"key,omitempty"`
	Data ListTileWithSubheaderTitle `json:"data,omitempty"`
}
