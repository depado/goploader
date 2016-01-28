package utils

import (
	"html/template"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
)

// InitTemplates loads the HTML templates from the rice box
func InitTemplates(r *gin.Engine, tbox *rice.Box, names ...string) error {
	var err error
	var tmpl string
	var message *template.Template
	for _, x := range names {
		if tmpl, err = tbox.String(x); err != nil {
			return err
		}
		if message, err = template.New(x).Parse(tmpl); err != nil {
			return err
		}
		r.SetHTMLTemplate(message)
	}
	return err
}
