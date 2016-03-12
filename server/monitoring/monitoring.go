package monitoring

import (
	"time"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/models"
)

// Monit monitors the database and file system to remove old entries
func Monit() {
	logger.Info("monitoring", "Started Monitoring")
	var err error
	tc := time.NewTicker(1 * time.Minute)
	for {
		err = database.DB.View(func(tx *bolt.Tx) error {
			now := time.Now()
			found := 0
			b := tx.Bucket([]byte("resources"))
			b.ForEach(func(k, v []byte) error {
				r := &models.Resource{}
				err := r.Decode(v)
				if err != nil {
					return err
				}
				if r.DeleteAt.Before(now) {
					found++
					r.Delete()
				}
				return nil
			})
			if found > 0 {
				logger.Info("monitoring", "Cleaning DB", "Flushed", found, "DB entries and files")
			}
			return nil
		})
		if err != nil {
			logger.Err("monitoring", "While checking", err)
		}
		<-tc.C
	}
}
