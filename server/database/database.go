package database

import (
	"time"

	"github.com/asdine/storm"
	bolt "go.etcd.io/bbolt"

	"github.com/Depado/goploader/server/conf"
)

// DB is the main database. Put in separate package for use in external ones.
var DB *storm.DB

// Initialize initializes the database (creating it if necessary)
func Initialize() error {
	var err error
	DB, err = storm.Open(conf.C.DB, storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	return err
}
