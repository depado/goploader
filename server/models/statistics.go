package models

import (
	"fmt"

	"github.com/asdine/storm"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/utils"
)

// Statistics is the struct representing the server statistics
type Statistics struct {
	ID           int
	TotalSize    uint64 `json:"total_size"`
	TotalFiles   uint64 `json:"total_files"`
	CurrentSize  uint64 `json:"current_size"`
	CurrentFiles uint64 `json:"current_files"`
}

// S is the exported main statistics structure
var S Statistics

// Initialize initializes the buckets if necessary
func Initialize() error {
	var err error
	var cfiles uint64
	var csize uint64
	var rs []Resource

	if err = database.DB.Init(&Resource{}); err != nil {
		logger.Fatal("server", "Couldn't initialize bucket", err)
	}
	if err = database.DB.Init(&Statistics{}); err != nil {
		logger.Fatal("server", "Couldn't initialize bucket", err)
	}
	logger.Debug("server", "Started Initialize on statistics object")
	if err = database.DB.One("ID", 0, &S); err != nil {
		if err == storm.ErrNotFound {
			if err = database.DB.All(&rs); err != nil {
				if err == storm.ErrNotFound {
					return nil
				}
				return err
			}
			for _, r := range rs {
				S.CurrentFiles++
				S.CurrentSize += uint64(r.Size)
			}
		} else {
			return err
		}
	}
	logger.Info("server", fmt.Sprintf("Total   %d (%s)", S.TotalFiles, utils.HumanBytes(S.TotalSize)))
	logger.Info("server", fmt.Sprintf("Current %d (%s)", cfiles, utils.HumanBytes(csize)))
	logger.Debug("server", "Done Initialize on statistics object")
	return err
}
