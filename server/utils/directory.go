package utils

import (
	"log"
	"os"
)

// EnsureDir checks if a folder exists, if not, creates it.
func EnsureDir(fp string) error {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("[INFO][System]\tCreating %v directory.\n", fp)
		os.Mkdir(fp, 0777)
	} else if err != nil {
		return err
	}
	return nil
}
