package types

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

type Schedule struct {
	App          string       `json:"app"`
	Cycle        string       `json:"cycle"`
	Datasource   string       `json:"datasource"`
	ID           cron.EntryID `json:"id"`
	Organization string       `json:"organization"`
	Schedule     string       `json:"schedule"`
	Type         string       `json:"type"`
	UpdatedAt    time.Time    `json:"updatedAt,omitempty"`
}
