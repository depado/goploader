package setup

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/utils"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "setup.html", gin.H{})
}

func configure(c *gin.Context) {
	var form conf.Conf
	var dat []byte
	var err error

	if err = c.Bind(&form); err == nil {
		errors := form.Validate()
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, errors)
			return
		}
		if err = form.FillDefaults(); err != nil {
			fmt.Println("An error occured while filling default values :", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if dat, err = yaml.Marshal(&form); err != nil {
			fmt.Println("An error occured while marshalling the yaml data :", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if err = ioutil.WriteFile("conf.yml", dat, 0644); err != nil {
			fmt.Println("An error occured while writing the conf.yml file :", err)
			c.AbortWithError(http.StatusInternalServerError, err)
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
	tbox, _ := rice.FindBox("templates")
	abox, _ := rice.FindBox("assets")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if err = utils.InitAssetsTemplates(r, tbox, abox, false, "setup.html"); err != nil {
		log.Fatal(err)
	}
	r.GET("/", index)
	r.POST("/", configure)
	fmt.Println("Please go to http://127.0.0.1:8008 to setup goploader.")
	r.Run(":8008")
}
