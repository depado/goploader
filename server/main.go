package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	flag "github.com/ogier/pflag"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/monitoring"
	"github.com/Depado/goploader/server/setup"
	"github.com/Depado/goploader/server/utils"
	"github.com/Depado/goploader/server/views"
)

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
	if database.DB, err = gorm.Open("sqlite3", conf.C.DB); err != nil {
		log.Fatal(err)
	}
	database.DB.AutoMigrate(&models.ResourceEntry{})

	go monitoring.Monit(&database.DB)

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
	r.GET("/", views.Index)
	r.POST("/", views.Create)
	r.GET("/v/:uniuri/:key", views.View)
	// Run
	r.Run(fmt.Sprintf(":%d", conf.C.Port))
}
