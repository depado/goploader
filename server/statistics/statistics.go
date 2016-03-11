package statistics

import (
	"encoding/json"

	"github.com/Depado/goploader/server/database"
	"github.com/boltdb/bolt"
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
	var err error
	var data []byte

	if data, err = s.Encode(); err != nil {
		return err
	}
	return database.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("statistics")).Put([]byte("main"), data)
	})
}

// Initialize loads the previous state of the statistics
func Initialize() {
	sp := &S
	database.DB.View(func(tx *bolt.Tx) error {
		return sp.Decode(tx.Bucket([]byte("statistics")).Get([]byte("main")))
	})
}

// Encode encodes a Resource to JSON
func (s Statistics) Encode() ([]byte, error) {
	return json.Marshal(s)
}

// Decode decodes a JSON struct to Resource
func (s *Statistics) Decode(data []byte) error {
	return json.Unmarshal(data, s)
}
