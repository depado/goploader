package models

import (
	"encoding/json"
	"fmt"

	"github.com/asdine/storm"
	bolt "go.etcd.io/bbolt"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
)

// Initialize initializes the buckets and
func Initialize() error {
	var err error

	if err = database.DB.Init(&Resource{}); err != nil {
		logger.Fatal("server", "Couldn't initialize bucket", err)
	}
	Migrate()
	if err = S.Load(); err != nil {
		switch err {
		case storm.ErrNotFound:
			if err = S.Evaluate(); err != nil {
				return err
			}
		default:
			return err
		}
	}

	S.Info()
	logger.Debug("server", "Done Initialize on statistics object")
	return err
}

// Migrate migrates the old data format to the new one
func Migrate() {
	var err error
	type us struct {
		TotalSize    uint64
		TotalFiles   uint64
		CurrentSize  uint64
		CurrentFiles uint64
	}
	var ss us
	if err = database.DB.Get("statistics", "main", &ss); err == nil {
		S.CurrentFiles = ss.CurrentFiles
		S.CurrentSize = ss.CurrentSize
		S.TotalFiles = ss.TotalFiles
		S.TotalSize = ss.TotalSize
		if err = S.Save(); err != nil {
			logger.Err("server", "Couldn't migrate old stats", err)
			return
		}
		logger.Info("server", "Had to migrate old stats")
		database.DB.Delete("statistics", "main")
		database.DB.Drop("statistics")
	}

	var rr []Resource
	var r Resource
	err = database.DB.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("no bucket")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err = json.Unmarshal(v, &r); err == nil {
				r.UnixDeleteAt = r.DeleteAt.Unix()
				rr = append(rr, r)
			}
		}
		return nil
	})
	if err == nil {
		for _, r := range rr {
			database.DB.Save(&r)
		}
		logger.Info("server", "Had to migrate old resources")
		database.DB.Drop("resources")
	}
}
