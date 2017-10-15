package models

import (
	"fmt"

	"github.com/asdine/storm"

	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/utils"
)

const (
	statsBucket = "stats"
	statsKey    = "main"
)

// S is the exported main statistics structure
var S Statistics

// Statistics is the struct representing the server statistics
type Statistics struct {
	TotalSize    uint64 `json:"total_size"`
	TotalFiles   uint64 `json:"total_files"`
	CurrentSize  uint64 `json:"current_size"`
	CurrentFiles uint64 `json:"current_files"`
}

// Load loads the statistics from the database
func (s *Statistics) Load() error {
	return database.DB.Get(statsBucket, statsKey, s)
}

// Save saves the statistics to the database
func (s *Statistics) Save() error {
	return database.DB.Set(statsBucket, statsKey, s)
}

// Add adds the given resource to the statistics and saves it to the database
func (s *Statistics) Add(r Resource) error {
	s.TotalFiles++
	s.CurrentFiles++
	s.TotalSize += uint64(r.Size)
	s.CurrentSize += uint64(r.Size)

	return s.Save()
}

// Remove removes the given resource from the current statistics and saves it
// to the database
func (s *Statistics) Remove(r Resource) error {
	s.CurrentFiles--
	s.CurrentSize -= uint64(r.Size)

	return s.Save()
}

// Evaluate recalculates the statistics from what's present in the database and
// save the statistics to the database
func (s *Statistics) Evaluate() error {
	var err error
	var rs []Resource

	// Find all the resources in the database and recalculate the stats
	if err = database.DB.All(&rs); err != nil {
		if err == storm.ErrNotFound {
			return nil
		}
		return err
	}
	// For all the present files update the statistics object
	for _, r := range rs {
		s.CurrentFiles++
		s.CurrentSize += uint64(r.Size)
		s.TotalFiles++
		s.TotalSize += uint64(r.Size)
	}
	return s.Save()
}

// Info prints out the information about statistics in the logs
func (s Statistics) Info() {
	logger.Info("server", fmt.Sprintf("Total   %d (%s)", s.TotalFiles, utils.HumanBytes(s.TotalSize)))
	logger.Info("server", fmt.Sprintf("Current %d (%s)", s.CurrentFiles, utils.HumanBytes(s.CurrentSize)))
}
