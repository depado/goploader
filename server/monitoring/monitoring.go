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
	"github.com/Depado/goploader/server/utils"
)

// Monit monitors the database and file system to remove old entries
func Monit() {
	logger.Info("monitoring", "Started Monitoring")
	var err error
	tc := time.NewTicker(1 * time.Minute)
	for {
		logger.Debug("monitoring", "Started Monit on Resources")
		now := time.Now()
		found := 0
		var sizeRemoved uint64
		err = database.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("resources"))
			err = b.ForEach(func(k, v []byte) error {
				var err error
				r := &models.Resource{}
				if err = r.Decode(v); err != nil {
					return err
				}
				if r.DeleteAt.Before(now) {
					if err = b.Delete([]byte(r.Key)); err != nil {
						logger.Err("monitoring", "While deleting entry (skipped)", err)
					} else {
						if err = os.Remove(path.Join(conf.C.UploadDir, r.Key)); err != nil {
							logger.Err("monitoring", "While deleting file (skipped)", err)
						} else {
							found++
							sizeRemoved += uint64(r.Size)
						}
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
			if found > 0 {
				var data []byte
				models.S.CurrentFiles -= uint64(found)
				models.S.CurrentSize -= sizeRemoved
				if data, err = models.S.Encode(); err != nil {
					return err
				}
				return tx.Bucket([]byte("statistics")).Put([]byte("main"), data)
			}
			return nil
		})
		if err != nil {
			logger.Err("monitoring", "While monitoring", err)
		} else {
			if found > 0 {
				logger.Info("monitoring", fmt.Sprintf("Deleted %d entries and files in %s", found, time.Since(now)))
				logger.Info("monitoring", fmt.Sprintf("Serving %d (%s) files", models.S.CurrentFiles, utils.HumanBytes(models.S.CurrentSize)))
			}
			logger.Debug("monitoring", fmt.Sprintf("Serving %d (%s) files", models.S.CurrentFiles, utils.HumanBytes(models.S.CurrentSize)))
		}
		logger.Debug("monitoring", fmt.Sprintf("Done Monit on Resources (%s)", time.Since(now)))
		<-tc.C
	}
}
