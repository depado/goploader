package views

import (
	"log"
	"net/http"

	"github.com/Depado/goploader/server/conf"
	"github.com/gin-gonic/gin"
)

// Index handles the main page
func Index(c *gin.Context) {
	log.Printf("[INFO][%s]\tIssued a GET request\n", c.ClientIP())
	c.HTML(http.StatusOK, "index.html", gin.H{
		"fulldoc":    conf.C.FullDoc,
		"size_limit": conf.C.SizeLimit,
	})
}
