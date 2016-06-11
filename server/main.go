package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
	flag "github.com/ogier/pflag"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/database"
	"github.com/Depado/goploader/server/logger"
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

	tbox, _ := rice.FindBox("templates")
	abox, _ := rice.FindBox("assets")

	flag.StringVarP(&cp, "conf", "c", "conf.yml", "Local path to configuration file.")
	flag.BoolVarP(&initial, "initial", "i", false, "Run the initial setup of the server.")
	flag.Parse()

	if err = conf.Load(cp, !initial); err != nil || initial {
		setup.Run()
	}
	if err = database.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()
	if err = models.Initialize(); err != nil {
		log.Fatal(err)
	}
	go monitoring.Monit()

	logger.Info("server", "Started goploader server on port", conf.C.Port)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if !conf.C.NoWeb {
		if err = utils.InitAssetsTemplates(r, tbox, abox, true, "index.html"); err != nil {
			log.Fatal(err)
		}
		r.Static("/releases", "releases")
		r.GET("/", views.Index)
	}
	r.POST("/", views.Create)
	r.GET("/v/:uniuri/:key", views.View)
	r.HEAD("/v/:uniuri/:key", views.Head)
	if conf.C.ServeHTTPS {
		http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.C.Port), conf.C.SSLCert, conf.C.SSLPrivKey, r)
	} else {
		if err = r.Run(fmt.Sprintf(":%d", conf.C.Port)); err != nil {
			logger.Err("server", "Fatal error", err)
		}
	}
}
