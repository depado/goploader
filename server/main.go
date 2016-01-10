package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nu7hatch/gouuid"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/monitoring"
	"github.com/Depado/goploader/server/utils"
)

var db gorm.DB

func index(c *gin.Context) {
	log.Printf("[INFO][%s]\tIssued a GET request\n", c.ClientIP())
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func create(c *gin.Context) {
	var err error
	var u *uuid.UUID
	remote := c.ClientIP()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 20000000)

	fd, h, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring reading file : %s", remote, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer fd.Close()
	if u, err = uuid.NewV4(); err != nil {
		log.Printf("[ERROR][%s]\tDuring creation of uuid : %s\n", remote, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	path := path.Join(conf.C.UploadDir, u.String())
	file, err := os.Create(path)
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring file creation : %s\n", remote, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	wr, err := io.Copy(file, bufio.NewReaderSize(fd, 512))
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring writing file : %s\n", remote, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	db.Create(&models.ResourceEntry{Key: u.String(), Name: h.Filename})
	log.Printf("[INFO][%s]\tCreated %s file and entry (%v bytes written)\n", remote, u.String(), wr)
	c.String(http.StatusCreated, "https://"+conf.C.NameServer+"/v/"+u.String()+"\n")
}

func view(c *gin.Context) {
	id := c.Param("uuid")
	re := models.ResourceEntry{}
	remote := c.ClientIP()

	db.Where(&models.ResourceEntry{Key: id}).First(&re)
	if re.Key == "" {
		log.Printf("[INFO][%s]\tNot found : %s", remote, id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Printf("[INFO][%s]\tFetched %s file and entry\n", remote, id)
	f, err := os.Open(conf.C.UploadDir + re.Key)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Disposition", "filename=\""+re.Name+"\"")
	http.ServeContent(c.Writer, c.Request, re.Key, re.CreatedAt, f)
}

func main() {
	var err error

	confPath := flag.String("c", "conf.yml", "Local path to configuration file.")
	flag.Parse()

	if err = conf.Load(*confPath); err != nil {
		log.Fatal(err)
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
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")
	r.Static("/favicon.ico", "./assets/favicon.ico")
	r.GET("/", index)
	r.POST("/", create)
	r.GET("/v/:uuid", view)
	r.Run(fmt.Sprintf(":%d", conf.C.Port))
}
