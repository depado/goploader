package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/utils"
	"github.com/Depado/goploader/server/views"
)

// Setup creates the gin Engine
func Setup() (*gin.Engine, error) {
	var err error

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if !conf.C.NoWeb {
		if err = utils.InitAssetsTemplates(r, true, "index.html", "mobile.html", "code.tmpl"); err != nil {
			return nil, err
		}
		r.Static("/releases", "releases")
		r.GET("/", views.Index)
		r.GET("/sw.js", func(c *gin.Context) {
			c.File("sw.js")
		})
		r.GET("/simple", views.SimpleIndex)

	}
	if conf.C.DisableEncryption {
		r.POST("/", views.Create)
		r.GET("/v/:uniuri", views.View)
		r.HEAD("/v/:uniuri", views.Head)
		r.GET("/v/:uniuri/:lang", views.ViewCode)
	} else {
		r.POST("/", views.CreateC)
		r.GET("/v/:uniuri/:key", views.ViewC)
		r.HEAD("/v/:uniuri/:key", views.HeadC)
		r.GET("/v/:uniuri/:key/:lang", views.ViewCCode)
		//TODO: Add HEAD for ViewCCode
	}
	if conf.C.PrometheusEnabled {
		r.Any("/metrics", gin.WrapH(promhttp.Handler()))
	}
	return r, nil
}
