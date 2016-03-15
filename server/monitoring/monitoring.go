package monitoring

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/conf"
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
		if conf.C.Debug {
			logger.Info("monitoring", "Started Monit on Resources")
		}
		found := 0
		err = database.DB.Update(func(tx *bolt.Tx) error {
			now := time.Now()
			b := tx.Bucket([]byte("resources"))
			return b.ForEach(func(k, v []byte) error {
				var err error
				r := &models.Resource{}
				if err = r.Decode(v); err != nil {
					return err
				}
				if r.DeleteAt.Before(now) {
					found++
					if err = b.Delete([]byte(r.Key)); err != nil {
						return err
					}
					if err = os.Remove(path.Join(conf.C.UploadDir, r.Key)); err != nil {
						return err
					}
				}
				return nil
			})
		})
		if err != nil {
			logger.Err("monitoring", "While monitoring", err)
		} else {
			if found > 0 {
				logger.Info("monitoring", fmt.Sprintf("Deleted %d entries and files", found))
			}
		}
		if conf.C.Debug {
			logger.Info("monitoring", "Done Monit on Resources")
		}
		<-tc.C
	}
}
