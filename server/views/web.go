package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/statistics"
	"github.com/Depado/goploader/server/utils"
)

// Index handles the main page
func Index(c *gin.Context) {
	logger.InfoC(c, "server", "GET /")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"fulldoc":     conf.C.FullDoc,
		"size_limit":  conf.C.SizeLimit,
		"total_size":  utils.HumanBytes(statistics.S.TotalSize),
		"total_files": statistics.S.TotalFiles,
	})
}
