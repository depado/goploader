package monitoring

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/models"
)

var timeLimit = 2 * time.Hour

// Monit monitors the database and file system to remove old entries
func Monit(db *gorm.DB) {
	log.Println("[INFO][System]\tStarted monitoring of files and db entries")
	tc := time.NewTicker(1 * time.Minute)
	for {
		res := []models.ResourceEntry{}
		db.Find(&res, "created_at < ?", time.Now().Add(-timeLimit))
		db.Unscoped().Where("created_at < ?", time.Now().Add(-timeLimit)).Delete(&models.ResourceEntry{})
		if len(res) > 0 {
			log.Printf("[INFO][System]\tFlushing %d DB entries and files.\n", len(res))
		}
		for _, re := range res {
			err := os.Remove(path.Join(conf.C.UploadDir, re.Key))
			if err != nil {
				log.Printf("[ERROR][System]\tWhile deleting : %v", err)
			}
		}
		<-tc.C
	}
}
