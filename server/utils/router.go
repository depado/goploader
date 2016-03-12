package utils

import (
	"html/template"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/logger"
)

// InitAssetsTemplates initializes the router to use either a ricebox or the
// filesystem in case the ricebox couldn't be found.
func InitAssetsTemplates(r *gin.Engine, tbox, abox *rice.Box, verbose bool, names ...string) error {
	var err error

	if tbox != nil {
		mt := multitemplate.New()
		var tmpl string
		var message *template.Template
		for _, x := range names {
			if tmpl, err = tbox.String(x); err != nil {
				return err
			}
			if message, err = template.New(x).Parse(tmpl); err != nil {
				return err
			}
			mt.Add(x, message)
		}
		if verbose {
			logger.Info("server", "Loaded templates from \"templates\" box")
		}
		r.HTMLRender = mt
	} else {
		r.LoadHTMLGlob("templates/*")
		if verbose {
			logger.Info("server", "Loaded templates from disk")
		}
	}

	if abox != nil {
		r.StaticFS("/static", abox.HTTPBox())
		if verbose {
			logger.Info("server", "Loaded assets from \"assets\" box")
		}
	} else {
		r.Static("/static", "assets")
		if verbose {
			logger.Info("server", "Loaded assets from disk")
		}
	}
	return nil
}
