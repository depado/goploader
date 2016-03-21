package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/utils"
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
	logger.Debug("server", "Started Save on Resource", r.Key)
	var err error
	var data []byte

	if data, err = r.Encode(); err != nil {
		return err
	}
	err = database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resources")).Put([]byte(r.Key), data)
	})
	logger.Debug("server", "Done Save on Resource", r.Key)
	return err
}

// Get retrives the Resource from the bucket
func (r *Resource) Get(key string) error {
	logger.Debug("server", "Started Get on Resource", key)
	err := database.DB.View(func(tx *bolt.Tx) error {
		return r.Decode(tx.Bucket([]byte("resources")).Get([]byte(key)))
	})
	logger.Debug("server", "Done Get on Resource", r.Key)
	return err
}

// Delete deletes a resource in database and on disk
func (r Resource) Delete() error {
	logger.Debug("server", "Started Delete on Resource", r.Key)
	var err error
	err = database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("resources")).Delete([]byte(r.Key))
	})
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(conf.C.UploadDir, r.Key))
	logger.Debug("server", "Done Delete on Resource", r.Key)
	return err
}

// LogCreated logs when a file is created
func (r Resource) LogCreated(c *gin.Context) {
	e := fmt.Sprintf("%sCreated%s %s - %s", logger.Green, logger.Reset, r.Key, utils.HumanBytes(uint64(r.Size)))
	if r.Once {
		e += " - once"
	}
	logger.InfoC(c, "server", e)
}

// LogFetched logs when a file is fetched
func (r Resource) LogFetched(c *gin.Context) {
	e := fmt.Sprintf("%sFetched%s %s - %s", logger.Yellow, logger.Reset, r.Key, utils.HumanBytes(uint64(r.Size)))
	if r.Once {
		e += " - once"
	}
	logger.InfoC(c, "server", e)
}

// LogDeleted logs when a file is deleted (due to a one-time view)
func (r Resource) LogDeleted(c *gin.Context) {
	e := fmt.Sprintf("%sDeleted%s %s - %s", logger.Red, logger.Reset, r.Key, utils.HumanBytes(uint64(r.Size)))
	if r.Once {
		e += " - once"
	}
	logger.InfoC(c, "server", e)
}

// Encode encodes a Resource to JSON
func (r Resource) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode decodes a JSON struct to Resource
func (r *Resource) Decode(data []byte) error {
	return json.Unmarshal(data, r)
}
