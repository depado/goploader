package views

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/utils"
)

// Create handles the multipart form upload and creates a file
func Create(c *gin.Context) {
	var err error
	var duration time.Duration
	var once bool

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, conf.C.SizeLimit*utils.MegaByte)

	once = c.PostForm("once") != ""
	d := c.DefaultPostForm("duration", "1d")

	if val, ok := models.DurationMap[d]; ok {
		duration = val
	} else {
		logger.ErrC(c, "server", "Invalid duration", d)
		c.String(http.StatusBadRequest, "Invalid duration\n")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fd, h, err := c.Request.FormFile("file")
	if err != nil {
		logger.ErrC(c, "server", "Couldn't read file", err)
		c.String(http.StatusRequestEntityTooLarge, "Entity is too large (Max : %v MB)\n", conf.C.SizeLimit)
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
	defer fd.Close()

	u := uniuri.NewLen(conf.C.UniURILength)
	file, err := os.Create(path.Join(conf.C.UploadDir, u))
	if err != nil {
		logger.ErrC(c, "server", "Couldn't create file", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	wr, err := io.Copy(file, bufio.NewReaderSize(fd, 512))
	if err != nil {
		logger.ErrC(c, "server", "Couldn't write file", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if conf.C.DiskQuota > 0 {
		if models.S.CurrentSize+uint64(wr) > uint64(conf.C.DiskQuota*utils.GigaByte) {
			logger.ErrC(c, "server", "Quota exceeded")
			c.String(http.StatusBadRequest, "Not enough free space. Try again later.")
			c.AbortWithStatus(http.StatusBadRequest)
			os.Remove(path.Join(conf.C.UploadDir, u))
			return
		}
	}
	newres := &models.Resource{
		Key:      u,
		Name:     h.Filename,
		Once:     once,
		DeleteAt: time.Now().Add(duration),
		Size:     wr,
	}
	if err = newres.Save(); err != nil {
		logger.ErrC(c, "server", "Couldn't save in database", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	newres.LogCreated(c)
	ns := conf.C.NameServer
	if conf.C.AppendPort {
		ns = fmt.Sprintf("%s:%d", conf.C.NameServer, conf.C.Port)
	}
	c.String(http.StatusCreated, "%v://%s/v/%s\n", utils.DetectScheme(c), ns, u)
}

// View handles the file views
func View(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	re := models.Resource{}

	if err = re.Get(id); err != nil || re.Key == "" {
		logger.InfoC(c, "server", "Not found", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	re.LogFetched(c)
	f, err := os.Open(path.Join(conf.C.UploadDir, re.Key))
	if err != nil {
		logger.ErrC(c, "server", fmt.Sprintf("Couldn't open %s", re.Key), err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if conf.C.AlwaysDownload {
		c.Header("Content-Type", "application/octet-stream")
	}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	io.Copy(c.Writer, f)
	if re.Once {
		re.Delete()
		re.LogDeleted(c)
	}
}

// Head handles the head request for an encryptd file
func Head(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	re := models.Resource{}

	if err = re.Get(id); err != nil || re.Key == "" {
		logger.InfoC(c, "server", "Not found", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	re.LogFetched(c)
	f, err := os.Open(path.Join(conf.C.UploadDir, re.Key))
	if err != nil {
		logger.ErrC(c, "server", fmt.Sprintf("Couldn't open %s", re.Key), err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if conf.C.AlwaysDownload {
		c.Header("Content-Type", "application/octet-stream")
	}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	io.Copy(c.Writer, f)
}
