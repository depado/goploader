package models

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
)

// DurationMap is a map linking the received string and a time.Duration
var DurationMap = map[string]time.Duration{
	"30m": 30 * time.Minute,
	"1h":  1 * time.Hour,
	"6h":  6 * time.Hour,
	"1d":  24 * time.Hour,
	"1w":  24 * time.Hour * 7,
}

// Resource represents the data stored in the database
type Resource struct {
	Key      string
	Name     string
	Once     bool
	Size     int64
	DeleteAt time.Time
}

// Save writes the Resource to the bucket
func (r Resource) Save() error {
	var err error
	var data []byte

	if data, err = r.Encode(); err != nil {
		return err
	}
	return database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resources")).Put([]byte(r.Key), data)
	})
}

// Get retrives the Resource from the bucket
func (r *Resource) Get(key string) error {
	return database.DB.View(func(tx *bolt.Tx) error {
		return r.Decode(tx.Bucket([]byte("resources")).Get([]byte(key)))
	})
}

// Delete deletes a resource in database and on disk
func (r Resource) Delete() {
	var err error
	database.DB.Update(func(tx *bolt.Tx) error {
		if err = tx.Bucket([]byte("resources")).Delete([]byte(r.Key)); err != nil {
			log.Printf("[ERROR][System]\tWhile deleting : %v", err)
		}
		return nil
	})
	if err = os.Remove(path.Join(conf.C.UploadDir, r.Key)); err != nil {
		log.Printf("[ERROR][System]\tWhile deleting : %v", err)
	}
}

// Encode encodes a Resource to JSON
func (r Resource) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode decodes a JSON struct to Resource
func (r *Resource) Decode(data []byte) error {
	return json.Unmarshal(data, r)
}
