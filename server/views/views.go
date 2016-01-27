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

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/models"
	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

// Index handles the main page
func Index(c *gin.Context) {
	log.Printf("[INFO][%s]\tIssued a GET request\n", c.ClientIP())
	if conf.C.FullDoc {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	} else {
		c.HTML(http.StatusOK, "welcome.html", gin.H{})
	}
}

// Create handles the multipart form upload
func Create(c *gin.Context) {
	var err error
	remote := c.ClientIP()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, conf.C.LimitSize*1000000)

	fd, h, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring reading file : %s", remote, err)
		c.String(http.StatusRequestEntityTooLarge, "Entity is too large (Max : %v MB)\n", conf.C.LimitSize)
		c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
	defer fd.Close()

	u := uniuri.NewLen(conf.C.UniURILength)

	k := uniuri.NewLen(16)
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
	}
	database.DB.Create(&models.ResourceEntry{Key: u, Name: h.Filename})

	log.Printf("[INFO][%s]\tCreated %s file and entry (%v bytes written)\n", remote, u, wr)
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	c.String(http.StatusCreated, "%s://%s/v/%s/%s\n", scheme, conf.C.NameServer, u, k)
}

// View handles the file views
func View(c *gin.Context) {
	id := c.Param("uniuri")
	key := c.Param("key")
	re := models.ResourceEntry{}
	remote := c.ClientIP()

	database.DB.Where(&models.ResourceEntry{Key: id}).First(&re)
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
}
