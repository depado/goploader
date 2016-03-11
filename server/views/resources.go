package views

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/utils"
)

// Create handles the multipart form upload
func Create(c *gin.Context) {
	var err error
	var duration time.Duration
	var once bool

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, conf.C.SizeLimit*1000000)

	remote := c.ClientIP()
	once = c.PostForm("once") != ""
	d := c.DefaultPostForm("duration", "1d")

	if val, ok := models.DurationMap[d]; ok {
		duration = val
	} else {
		log.Printf("[ERROR][%s]\tInvalid duration : %s", remote, d)
		c.String(http.StatusBadRequest, "Invalid duration\n")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fd, h, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring reading file : %s", remote, err)
		c.String(http.StatusRequestEntityTooLarge, "Entity is too large (Max : %v MB)\n", conf.C.SizeLimit)
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
	defer fd.Close()

	u := uniuri.NewLen(conf.C.UniURILength)
	k := uniuri.NewLen(conf.C.KeyLength)
	kb := []byte(k)
	block, err := aes.NewCipher(kb)
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring Cipher creation : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBEncrypter(block, iv[:])

	path := path.Join(conf.C.UploadDir, u)
	file, err := os.Create(path)
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring file creation : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := &cipher.StreamWriter{S: stream, W: file}
	// Copy the input file to the output file, encrypting on the fly.
	wr, err := io.Copy(writer, bufio.NewReaderSize(fd, 512))
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring writing : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	newres := &models.Resource{
		Key:      u,
		Name:     h.Filename,
		Once:     once,
		DeleteAt: time.Now().Add(duration),
		Size:     wr,
	}
	if err = newres.Save(); err != nil {
		log.Printf("[ERROR][%s]\tDuring saving : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO][%s]\tCreated %s file and entry (%v bytes written) (%s lifetime)\n", remote, u, wr, d)
	c.String(http.StatusCreated, "%v://%s/v/%s/%s\n", utils.DetectScheme(c), conf.C.NameServer, u, k)
}

// View handles the file views
func View(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	key := c.Param("key")
	re := models.Resource{}
	remote := c.ClientIP()

	if err = re.Get(id); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if re.Key == "" {
		log.Printf("[INFO][%s]\tNot found : %s", remote, id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Printf("[INFO][%s]\tFetched %s file and entry\n", remote, id)
	f, err := os.Open(path.Join(conf.C.UploadDir, re.Key))
	if err != nil {
		log.Printf("[ERROR][%s]\tWhile opening %s file\n", remote, id)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring Cipher creation : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBDecrypter(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	io.Copy(c.Writer, reader)
	if re.Once {
		re.Delete()
	}
}

// Head handles the head request for a file
func Head(c *gin.Context) {
	var err error

	id := c.Param("uniuri")
	key := c.Param("key")
	re := models.Resource{}
	remote := c.ClientIP()

	if err = re.Get(id); err != nil {
		log.Printf("[ERROR][%s]\t%s", remote, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if re.Key == "" {
		log.Printf("[INFO][%s]\tNot found : %s", remote, id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Printf("[INFO][%s]\tHead %s file and entry\n", remote, id)
	f, err := os.Open(path.Join(conf.C.UploadDir, re.Key))
	if err != nil {
		log.Printf("[ERROR][%s]\tWhile opening %s file\n", remote, id)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring Cipher creation : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBDecrypter(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: f}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	io.Copy(c.Writer, reader)
}
