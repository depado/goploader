package views

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"

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
	token := c.PostForm("token")

	if conf.C.Token != "" && conf.C.Token != token {
		logger.ErrC(c, "server", "Incorrect token")
		c.String(http.StatusUnauthorized, "Incorrect token\n")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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
	defer fd.Close() //nolint:errcheck

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
			os.Remove(path.Join(conf.C.UploadDir, res.Key)) //nolint:errcheck
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
	c.String(http.StatusCreated, "%s/v/%s/%s\n", utils.ServerURI(c), res.Key, k)
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
	stream := cipher.NewCTR(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	if conf.C.AlwaysDownload {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=\""+re.Name+"\"")
	} else {
		c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		logger.ErrC(c, "server", fmt.Sprintf("Couldn't copy %s", re.Key), err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if re.Once {
		if err := re.Delete(); err != nil {
			logger.ErrC(c, "server", fmt.Sprintf("Couldn't delete %s", re.Key), err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		re.LogDeleted(c)
	}
}

// ViewCCode allows to see the file with syntax highliting and extra options
func ViewCCode(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	key := c.Param("key")
	lang := c.Param("lang")
	theme := c.DefaultQuery("theme", "dark")
	lines := c.Query("lines") == "true"
	re := models.Resource{}

	if err = re.Get(id); err != nil || re.Key == "" {
		logger.InfoC(c, "server", "Not found", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if re.Size > conf.C.ViewLimit*utils.MegaByte {
		logger.InfoC(c, "server", fmt.Sprintf("Tried to view %s but it is too large (%s > %s)", re.Key, utils.HumanBytes(uint64(re.Size)), utils.HumanBytes(uint64(conf.C.ViewLimit*utils.MegaByte))))
		c.AbortWithStatus(http.StatusForbidden)
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
	stream := cipher.NewCTR(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		logger.ErrC(c, "server", "Couldn't read from file", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	bb := buf.Bytes()
	c.HTML(http.StatusOK, "code.tmpl", gin.H{
		"code":  string(bb),
		"lang":  lang,
		"theme": theme,
		"lines": lines,
		"name":  re.Name,
	})
	if re.Once {
		if err := re.Delete(); err != nil {
			logger.ErrC(c, "server", fmt.Sprintf("Couldn't delete %s", re.Key), err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
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
	stream := cipher.NewCTR(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	if conf.C.AlwaysDownload {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=\""+re.Name+"\"")
	} else {
		c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	}
	if _, err := io.Copy(c.Writer, reader); err != nil {
		logger.ErrC(c, "server", fmt.Sprintf("Couldn't copy %s", re.Key), err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
