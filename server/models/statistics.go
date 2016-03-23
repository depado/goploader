package models

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/utils"
)

// Statistics is the struct representing the server statistics
type Statistics struct {
	TotalSize    uint64
	TotalFiles   uint64
	CurrentSize  uint64
	CurrentFiles uint64
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
func Initialize() error {
	logger.Debug("server", "Started Initialize on statistics object")
	var cfiles uint64
	var csize uint64
	var err error

	err = database.DB.View(func(tx *bolt.Tx) error {
		S.Decode(tx.Bucket([]byte("statistics")).Get([]byte("main")))
		return tx.Bucket([]byte("resources")).ForEach(func(k, v []byte) error {
			r := &Resource{}
			if err = r.Decode(v); err != nil {
				return err
			}
			csize += uint64(r.Size)
			cfiles++
			return nil
		})
	})
	if err != nil {
		logger.Err("server", "Could not initialize statistics")
		return err
	}
	S.CurrentFiles = cfiles
	S.CurrentSize = csize
	logger.Info("server", fmt.Sprintf("Total   %d (%s)", S.TotalFiles, utils.HumanBytes(S.TotalSize)))
	logger.Info("server", fmt.Sprintf("Current %d (%s)", cfiles, utils.HumanBytes(csize)))
	logger.Debug("server", "Done Initialize on statistics object")
	return err
}

// Encode encodes a Resource to JSON
func (s Statistics) Encode() ([]byte, error) {
	return json.Marshal(s)
}

// Decode decodes a JSON struct to Resource
func (s *Statistics) Decode(data []byte) error {
	return json.Unmarshal(data, s)
}
