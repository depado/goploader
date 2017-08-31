package views

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sourcegraph/syntaxhighlight"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/utils"
)

// CreateC handles the multipart form upload and creates an encrypted file
func CreateC(c *gin.Context) {
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

	res := models.NewResourceFromForm(h, once, duration)
	k, err := res.WriteEncrypted(fd)
	if err != nil {
		logger.ErrC(c, "server", "Couldn't write file", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if conf.C.DiskQuota > 0 {
		if models.S.CurrentSize+uint64(res.Size) > uint64(conf.C.DiskQuota*utils.GigaByte) {
			logger.ErrC(c, "server", "Quota exceeded")
			c.String(http.StatusBadRequest, "Insufficient disk space. Try again later.")
			c.AbortWithStatus(http.StatusBadRequest)
			os.Remove(path.Join(conf.C.UploadDir, res.Key))
			return
		}
	}

	if err = res.Save(); err != nil {
		logger.ErrC(c, "server", "Couldn't save in the database", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res.OnCreated(c)
	ns := conf.C.NameServer
	if conf.C.AppendPort {
		ns = fmt.Sprintf("%s:%d", conf.C.NameServer, conf.C.Port)
	}
	c.String(http.StatusCreated, "%v://%s/v/%s/%s\n", utils.DetectScheme(c), ns, res.Key, k)
}

// ViewC handles the file views for encrypted files
func ViewC(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	key := c.Param("key")
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
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.ErrC(c, "server", "Couldn't create AES cipher", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBDecrypter(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	if conf.C.AlwaysDownload {
		c.Header("Content-Type", "application/octet-stream")
	}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	if _, ok := c.GetQuery("code"); ok && !conf.C.AlwaysDownload {
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		bb := buf.Bytes()
		if bb, err = syntaxhighlight.AsHTML(bb); err != nil {
			logger.ErrC(c, "server", "Couldn't parse syntax of file", err)
			c.String(http.StatusInternalServerError, "Something went wrong on the server. Try again later.")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "code.tmpl", gin.H{"code": template.HTML(bb)})
	} else {
		io.Copy(c.Writer, reader)
	}
	if re.Once {
		re.Delete()
		re.LogDeleted(c)
	}
}

// HeadC handles the head request for an encryptd file
func HeadC(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	key := c.Param("key")
	re := models.Resource{}

	if err = re.Get(id); err != nil {
		logger.InfoC(c, "server", "Not found", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if re.Key == "" {
		logger.InfoC(c, "server", "Not found", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	logger.InfoC(c, "server", "Head", fmt.Sprintf("%s - %s", re.Key, utils.HumanBytes(uint64(re.Size))))
	f, err := os.Open(path.Join(conf.C.UploadDir, re.Key))
	if err != nil {
		logger.ErrC(c, "server", fmt.Sprintf("Couldn't open %s", re.Key), err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.ErrC(c, "server", "Couldn't create AES cipher", err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBDecrypter(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	io.Copy(c.Writer, reader)
}
