package statistics

import (
	"encoding/json"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
)

// Statistics is the struct representing the server statistics
type Statistics struct {
	TotalSize  uint64
	TotalFiles uint64
}

// S is the exported main statistics structure
var S Statistics

// Save writes the Resource to the bucket
func (s Statistics) Save() error {
	logger.Debug("server", "Started Save on statistics object")
	var err error
	var data []byte

	if data, err = s.Encode(); err != nil {
		return err
	}
	err = database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("statistics")).Put([]byte("main"), data)
	})
	logger.Debug("server", "Done Save on statistics object")
	return err
}

// Initialize loads the previous state of the statistics
func Initialize() {
	logger.Debug("server", "Started Initialize on statistics object")
	sp := &S
	database.DB.View(func(tx *bolt.Tx) error {
		return sp.Decode(tx.Bucket([]byte("statistics")).Get([]byte("main")))
	})
	logger.Debug("server", "Done Initialize on statistics object")
}

// Encode encodes a Resource to JSON
func (s Statistics) Encode() ([]byte, error) {
	return json.Marshal(s)
}

// Decode decodes a JSON struct to Resource
func (s *Statistics) Decode(data []byte) error {
	return json.Unmarshal(data, s)
}
