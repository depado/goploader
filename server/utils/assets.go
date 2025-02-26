package utils

import (
	"embed"
	"html/template"
	"net/http"
	"io/fs"

	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/logger"
)

//go:embed all:assets
var assets embed.FS
//go:embed all:templates
var templates embed.FS

// InitAssetsTemplates initializes the router to use embedded assets
func InitAssetsTemplates(r *gin.Engine, verbose bool, names ...string) error {
	template := template.Must(template.ParseFS(templates, "templates/*"))
	r.SetHTMLTemplate(template)
	logger.Debug("server", "Loaded templates")

	assetsFp, _ := fs.Sub(assets, "assets")
	r.StaticFS("/static", http.FS(assetsFp))
	logger.Debug("server", "Loaded assets")
	return nil
}
