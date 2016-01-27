package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	flag "github.com/ogier/pflag"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/monitoring"
	"github.com/Depado/goploader/server/setup"
	"github.com/Depado/goploader/server/utils"
)

var db gorm.DB

func index(c *gin.Context) {
	log.Printf("[INFO][%s]\tIssued a GET request\n", c.ClientIP())
	if conf.C.FullDoc {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	} else {
		c.HTML(http.StatusOK, "welcome.html", gin.H{})
	}
}

func create(c *gin.Context) {
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
	// No encryption : wr, err := io.Copy(file, bufio.NewReaderSize(fd, 512))
	// Copy the input file to the output file, encrypting as we go.
	wr, err := io.Copy(writer, bufio.NewReaderSize(fd, 512))
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring writing : %s\n", remote, err)
		c.String(http.StatusInternalServerError, "Something went wrong on the server side. Try again later.")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	db.Create(&models.ResourceEntry{Key: u, Name: h.Filename})

	log.Printf("[INFO][%s]\tCreated %s file and entry (%v bytes written)\n", remote, u, wr)
	c.String(http.StatusCreated, "https://%s/v/%s/%s\n", conf.C.NameServer, u, k)
}

func view(c *gin.Context) {
	id := c.Param("uniuri")
	key := c.Param("key")
	re := models.ResourceEntry{}
	remote := c.ClientIP()

	db.Where(&models.ResourceEntry{Key: id}).First(&re)
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

func main() {
	var err error
	var cp string
	var initial bool
	var conferr error

	flag.StringVarP(&cp, "conf", "c", "conf.yml", "Local path to configuration file.")
	flag.BoolVarP(&initial, "initial", "i", false, "Run the initial setup of the server.")
	flag.Parse()

	conferr = conf.Load(cp)
	if conferr != nil || initial {
		setup.Run()
	}
	if err = utils.EnsureDir(conf.C.UploadDir); err != nil {
		log.Fatal(err)
	}
	if db, err = gorm.Open("sqlite3", conf.C.DB); err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.ResourceEntry{})

	go monitoring.Monit(&db)

	log.Printf("[INFO][System]\tStarted goploader server on port %d\n", conf.C.Port)
	if !conf.C.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	// Default router
	r := gin.Default()
	// Templates and static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")
	r.Static("/favicon.ico", "./assets/favicon.ico")
	// Routes
	r.GET("/", index)
	r.POST("/", create)
	r.GET("/v/:uniuri/:key", view)
	// Run
	r.Run(fmt.Sprintf(":%d", conf.C.Port))
}
