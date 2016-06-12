package router

import (
	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/utils"
	"github.com/Depado/goploader/server/views"
)

// Setup creates the gin Engine
func Setup(tbox, abox *rice.Box) (*gin.Engine, error) {
	var err error

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if !conf.C.NoWeb {
		if err = utils.InitAssetsTemplates(r, tbox, abox, true, "index.html"); err != nil {
			return nil, err
		}
		r.Static("/releases", "releases")
		r.GET("/", views.Index)
	}
	if conf.C.DisableEncryption {
		r.POST("/", views.Create)
		r.GET("/v/:uniuri", views.View)
		r.HEAD("/v/:uniuri", views.Head)
	} else {
		r.POST("/", views.CreateC)
		r.GET("/v/:uniuri/:key", views.ViewC)
		r.HEAD("/v/:uniuri/:key", views.HeadC)
	}
	return r, nil
}
