package database

import (
	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/conf"
)

// DB is the main database. Put in separate package for use in external ones.
var DB *bolt.DB

// Initialize initializes the database (creating it if necessary)
func Initialize() error {
	var err error
	if DB, err = bolt.Open(conf.C.DB, 0600, nil); err != nil {
		return err
	}
	err = DB.Update(func(tx *bolt.Tx) error {
		var err error
		if _, err = tx.CreateBucketIfNotExists([]byte("resources")); err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("statistics"))
		return err
	})
	return err
}
