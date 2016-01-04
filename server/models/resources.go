package models

import "github.com/jinzhu/gorm"

// ResourceEntry represents the data stored in the database
type ResourceEntry struct {
	gorm.Model
	Key  string
	Name string
}
