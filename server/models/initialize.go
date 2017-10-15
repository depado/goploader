package models

import (
	"github.com/asdine/storm"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
)

// Initialize initializes the buckets and
func Initialize() error {
	var err error

	if err = database.DB.Init(&Resource{}); err != nil {
		logger.Fatal("server", "Couldn't initialize bucket", err)
	}
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
