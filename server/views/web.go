package views

import (
	"net/http"

	"github.com/Depado/goploader/server/conf"
	"github.com/gin-gonic/gin"
)

// Index handles the main page
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"fulldoc":    conf.C.FullDoc,
		"size_limit": conf.C.SizeLimit,
	})
}
