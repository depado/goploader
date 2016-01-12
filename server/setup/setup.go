package setup

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/braintree/manners"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
)

var srv *manners.GracefulServer

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "setup.html", gin.H{})
}

func configure(c *gin.Context) {
	var form conf.Conf
	var err error

	if err = c.Bind(&form); err == nil {
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
		c.HTML(400, "setup.html", gin.H{"errors": err})
		return
	}
	c.HTML(201, "setup_ok.html", gin.H{})
}

// Run runs the setup server which is used to configure the application on the
// first run or when the -i/--initial option is used.
func Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")
	r.Static("/favicon.ico", "./assets/favicon.ico")
	r.GET("/", index)
	r.POST("/", configure)
	fmt.Println("Please go to http://127.0.0.1:8008 to setup goploader.")
	r.Run(":8008")
}
