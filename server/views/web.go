package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/utils"
)

// Index handles the main page
func Index(c *gin.Context) {
	logger.InfoC(c, "server", "GET /")
	data := gin.H{
		"fulldoc":        conf.C.FullDoc,
		"size_limit":     conf.C.SizeLimit,
		"sensitive_mode": conf.C.SensitiveMode,
	}
	if conf.C.Stats {
		data["total_size"] = utils.HumanBytes(models.S.TotalSize)
		data["total_files"] = models.S.TotalFiles
	}
	c.HTML(http.StatusOK, "index.html", data)
}
