package models

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
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
	if conf.C.Debug {
		logger.Info("server", "Started Save on Resource", r.Key)
	}
	var err error
	var data []byte

	if data, err = r.Encode(); err != nil {
		return err
	}
	err = database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resources")).Put([]byte(r.Key), data)
	})
	if conf.C.Debug {
		logger.Info("server", "Done Save on Resource", r.Key)
	}
	return err
}

// Get retrives the Resource from the bucket
func (r *Resource) Get(key string) error {
	if conf.C.Debug {
		logger.Info("server", "Started Get on Resource", key)
	}
	err := database.DB.View(func(tx *bolt.Tx) error {
		return r.Decode(tx.Bucket([]byte("resources")).Get([]byte(key)))
	})
	if conf.C.Debug {
		logger.Info("server", "Done Get on Resource", r.Key)
	}
	return err
}

// Delete deletes a resource in database and on disk
func (r Resource) Delete() error {
	var err error
	if conf.C.Debug {
		logger.Info("server", "Started Delete on Resource", r.Key)
	}
	if err = database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resources")).Delete([]byte(r.Key))
	}); err != nil {
		return err
	}
	if err = os.Remove(path.Join(conf.C.UploadDir, r.Key)); err != nil {
		return err
	}
	return nil
}

// Encode encodes a Resource to JSON
func (r Resource) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode decodes a JSON struct to Resource
func (r *Resource) Decode(data []byte) error {
	return json.Unmarshal(data, r)
}
