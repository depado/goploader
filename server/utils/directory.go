package utils

import (
	"os"

	"github.com/Depado/goploader/server/logger"
)

// EnsureDir checks if a folder exists, if not, creates it.
func EnsureDir(fp string) error {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		logger.Info("server", "Creating directory", fp)
		os.Mkdir(fp, 0777)
	} else if err != nil {
		return err
	}
	return nil
}
