package models

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/metrics"
	"github.com/Depado/goploader/server/utils"
)

// DurationMap is a map linking the received string and a time.Duration
var DurationMap = map[string]time.Duration{
	"30m": 30 * time.Minute,
	"1h":  1 * time.Hour,
	"6h":  6 * time.Hour,
	"1d":  24 * time.Hour,
	"1w":  24 * time.Hour * 7,
}

// Resource represents the data stored in the database
type Resource struct {
	Key      string    `json:"key" storm:"id,index"`
	Name     string    `json:"name"`
	Once     bool      `json:"once"`
	Size     int64     `json:"size"`
	DeleteAt time.Time `json:"delete_at" storm:"index"`
	Duration string    `json:"-"`
}

// NewResourceFromForm returns a new Resource instance with some fields calculated
func NewResourceFromForm(h *multipart.FileHeader, once bool, duration time.Duration) Resource {
	return Resource{
		Key:      uniuri.NewLen(conf.C.UniURILength),
		Name:     h.Filename,
		Once:     once,
		DeleteAt: time.Now().Add(duration),
	}
}

// NewStreamWriter creates a new encrypted AES stream writer with the given key
// and the given file descriptor
func (r Resource) NewStreamWriter(fd *os.File, key []byte) (*cipher.StreamWriter, error) {
	var block cipher.Block
	var err error

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBEncrypter(block, iv[:])
	return &cipher.StreamWriter{S: stream, W: fd}, nil
}

// Write isn't implemented yet
func (r Resource) Write() error {
	return nil
}

// WriteEncrypted is a method to write the file and encrypt it on the fly
// it returns the key that was used to encrypt the file so it can be sent back
// to the client.
func (r *Resource) WriteEncrypted(fd multipart.File) (string, error) {
	file, err := os.Create(path.Join(conf.C.UploadDir, r.Key))
	if err != nil {
		return "", err
	}
	defer file.Close()
	k := uniuri.NewLen(conf.C.KeyLength)
	kb := []byte(k)
	sw, err := r.NewStreamWriter(file, kb)
	if err != nil {
		return "", err
	}
	wr, err := io.Copy(sw, bufio.NewReaderSize(fd, 512))
	if err != nil {
		os.Remove(path.Join(conf.C.UploadDir, r.Key))
		return "", err
	}
	r.Size = wr
	return k, nil
}

// Save writes the Resource to the bucket
func (r Resource) Save() error {
	logger.Debug("server", "Started Save on Resource", r.Key)
	var err error
	// var data []byte

	if err = database.DB.Save(&r); err != nil {
		return err
	}

	// if data, err = r.Encode(); err != nil {
	// 	return err
	// }
	// err = database.DB.Update(func(tx *bolt.Tx) error {
	// 	if err = tx.Bucket([]byte("resources")).Put([]byte(r.Key), data); err != nil {
	// 		return err
	// 	}
	// 	S.TotalFiles++
	// 	S.TotalSize += uint64(r.Size)
	// 	S.CurrentFiles++
	// 	S.CurrentSize += uint64(r.Size)
	// 	if data, err = S.Encode(); err != nil {
	// 		return err
	// 	}
	// 	return tx.Bucket([]byte("statistics")).Put([]byte("main"), data)
	// })
	logger.Debug("server", "Done Save on Resource", r.Key)
	return err
}

// Get retrives the Resource from the bucket
func (r *Resource) Get(key string) error {
	logger.Debug("server", "Started Get on Resource", key)
	err := database.DB.One("Key", key, r)
	logger.Debug("server", "Done Get on Resource", key)
	return err
}

// Delete deletes a resource in database and on disk
func (r Resource) Delete() error {
	logger.Debug("server", "Started Delete on Resource", r.Key)
	var err error
	if err = database.DB.DeleteStruct(&r); err != nil {
		return err
	}
	S.CurrentFiles--
	S.CurrentSize -= uint64(r.Size)
	database.DB.Save(&S)
	err = os.Remove(path.Join(conf.C.UploadDir, r.Key))
	logger.Debug("server", "Done Delete on Resource", r.Key)
	logger.Debug("server", fmt.Sprintf("Serving %d (%s) files", S.CurrentFiles, utils.HumanBytes(S.CurrentSize)))
	return err
}

// LogCreated logs when a file is created
func (r Resource) LogCreated(c *gin.Context) {
	e := fmt.Sprintf("%sCreated%s %s - %s", logger.Green, logger.Reset, r.Key, utils.HumanBytes(uint64(r.Size)))
	if r.Once {
		e += " - once"
	}
	logger.InfoC(c, "server", e)
}

// LogFetched logs when a file is fetched
func (r Resource) LogFetched(c *gin.Context) {
	e := fmt.Sprintf("%sFetched%s %s - %s", logger.Yellow, logger.Reset, r.Key, utils.HumanBytes(uint64(r.Size)))
	if r.Once {
		e += " - once"
	}
	logger.InfoC(c, "server", e)
}

// LogDeleted logs when a file is deleted (due to a one-time view)
func (r Resource) LogDeleted(c *gin.Context) {
	e := fmt.Sprintf("%sDeleted%s %s - %s", logger.Red, logger.Reset, r.Key, utils.HumanBytes(uint64(r.Size)))
	if r.Once {
		e += " - once"
	}
	logger.InfoC(c, "server", e)
}

// OnCreated is the function called once the resource has been created and saved
func (r Resource) OnCreated(c *gin.Context) {
	r.LogCreated(c)
	if conf.C.PrometheusEnabled {
		r.logMetricsCreated(c)
	}
}

// LogMetricsCreated updates the prometheus metrics
func (r Resource) logMetricsCreated(c *gin.Context) {
	metrics.UploadedFilesSizeTotal.Add(float64(r.Size))
	metrics.UploadedFilesTotal.WithLabelValues(c.ClientIP()).Inc()
}
