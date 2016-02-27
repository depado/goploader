package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DurationMap is a map linking the received string and a time.Duration
var DurationMap = map[string]time.Duration{
	"30m": 30 * time.Minute,
	"1h":  1 * time.Hour,
	"6h":  6 * time.Hour,
	"1d":  24 * time.Hour,
	"1w":  168 * time.Hour,
}

// ResourceEntry represents the data stored in the database
type ResourceEntry struct {
	gorm.Model
	Key      string
	Name     string
	Once     bool
	DeleteAt time.Time
}
