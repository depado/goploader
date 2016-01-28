package setup

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/utils"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "setup.html", gin.H{})
}

func configure(c *gin.Context) {
	var form conf.UnparsedConf
	var err error

	if err = c.Bind(&form); err == nil {
		errors := form.Validate()
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, errors)
			return
		}
		d, err := yaml.Marshal(&form)
		if err != nil {
			fmt.Println("An error occured while marshalling the yaml data :", err)
			c.AbortWithError(500, err)
			return
		}
		if err = ioutil.WriteFile("conf.yml", d, 0644); err != nil {
			fmt.Println("An error occured while writing the conf.yml file :", err)
			c.AbortWithError(500, err)
			return
		}
	} else {
		fmt.Println("An error occured while reading the form data :", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// Run runs the setup server which is used to configure the application on the
// first run or when the -i/--initial option is used.
func Run() {
	var err error

	assetsBox, err := rice.FindBox("assets")
	if err != nil {
		log.Fatal(err)
	}
	templateBox, err := rice.FindBox("templates")
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	if err = utils.InitTemplates(r, templateBox, "setup.html"); err != nil {
		log.Fatal(err)
	}
	r.StaticFS("/static", assetsBox.HTTPBox())
	//r.Static("/favicon.ico", "./assets/favicon.ico")
	r.GET("/", index)
	r.POST("/", configure)
	fmt.Println("Please go to http://127.0.0.1:8008 to setup goploader.")
	r.Run(":8008")
}
