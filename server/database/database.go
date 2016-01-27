package database

import "github.com/jinzhu/gorm"

// DB is the main database. Put in separate package for use in external ones.
var DB gorm.DB
