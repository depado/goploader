package monitoring

import (
	"log"
	"time"

	"github.com/boltdb/bolt"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/models"
)

// Monit monitors the database and file system to remove old entries
func Monit() {
	log.Println("[INFO][System]\tStarted monitoring of files and db entries")
	tc := time.NewTicker(1 * time.Minute)
	for {
		database.DB.View(func(tx *bolt.Tx) error {
			now := time.Now()
			found := 0
			b := tx.Bucket([]byte("resources"))
			b.ForEach(func(k, v []byte) error {
				var err error
				r := &models.Resource{}
				if err = r.Decode(v); err != nil {
					return err
				}
				if r.DeleteAt.Before(now) {
					found++
					r.Delete()
				}
				return nil
			})
			if found > 0 {
				log.Printf("[INFO][System]\tFlushed %d DB entries and files.\n", found)
			}
			return nil
		})
		<-tc.C
	}
}
