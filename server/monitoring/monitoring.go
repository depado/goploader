package monitoring

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/asdine/storm"

	"github.com/asdine/storm/q"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/utils"
)

// Monit monitors the database and file system to remove old entries
func Monit() {
	logger.Info("monitoring", "Started Monitoring")

	// Execute the function when the server starts
	FindAndDelete()

	tc := time.NewTicker(1 * time.Minute)
	for {
		<-tc.C
		FindAndDelete()
	}
}

// FindAndDelete is a function to find and delete all outdated resources
func FindAndDelete() {
	var err error
	var tx storm.Node
	now := time.Now()

	var rs []models.Resource
	if err = database.DB.Select(q.Lt("UnixDeleteAt", now.Unix())).Find(&rs); err != nil {
		if err == storm.ErrNotFound {
			logger.Debug("monitoring", fmt.Sprintf("Done Monit on Resources (%s)", time.Since(now)))
			return
		}
		logger.Err("monitoring", "While monitoring", err)
		return
	}

	if tx, err = database.DB.Begin(true); err != nil {
		logger.Err("monitoring", "Couldn't start transaction", err)
	}
	for _, r := range rs {
		if err = os.Remove(path.Join(conf.C.UploadDir, r.Key)); err != nil {
			logger.Err("monitoring", "While deleting file (skipped)", err)
		}
		if err = tx.DeleteStruct(&r); err != nil {
			logger.Err("monitoring", "While deleting file from DB (skipped)", err)
		}
		models.S.CurrentSize -= uint64(r.Size)
		models.S.CurrentFiles--
	}
	if err = tx.Commit(); err != nil {
		logger.Err("monitoring", "Couldn't commit transaction", err)
	}
	if err = models.S.Save(); err != nil {
		logger.Err("monitoring", "Error while saving statistics", err)
	}
	logger.Info("monitoring", fmt.Sprintf("Deleted %d entries and files in %s", len(rs), time.Since(now)))
	logger.Info("monitoring", fmt.Sprintf("Serving %d (%s) files", models.S.CurrentFiles, utils.HumanBytes(models.S.CurrentSize)))
	logger.Debug("monitoring", fmt.Sprintf("Done Monit on Resources (%s)", time.Since(now)))
}
