package types

//Event is a struct containing a generic event created by any 9 Spokes service
type Event struct {
	Action    string                 `json:"action"`
	Service   string                 `json:"service"`
	Timestamp int64                  `json:"timestamp"`
	User      string                 `json:"user"`
	Email     string                 `json:"email"`
	Session   string                 `json:"session"`
	Data      map[string]interface{} `json:"data"`
	Channel   string                 `json:"channel"`
}

//HandlerConfig is entry in the handlers.yaml configuration file, used to supply runtime options to various event handlers
type HandlerConfig struct {
	Name    string                 `yaml:"name"`
	Enabled bool                   `yaml:"enabled"`
	Events  string                 `yaml:"events"`
	Data    map[string]interface{} `yaml:"data"`
}
