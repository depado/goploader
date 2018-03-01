package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"strings"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/curl"
	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/models"
	"github.com/Depado/goploader/server/utils"
)

// Index handles the main page
func Index(c *gin.Context) {
	logger.InfoC(c, "server", "GET /")
	if strings.HasPrefix(c.Request.Header.Get("User-Agent"), "curl") {
		curl.WriteTutorial(c)
		return
	}
	data := gin.H{
		"token":          conf.C.Token,
		"fulldoc":        conf.C.FullDoc,
		"size_limit":     utils.HumanBytes(uint64(conf.C.SizeLimit * utils.MegaByte)),
		"sensitive_mode": conf.C.SensitiveMode,
	}
	if conf.C.Stats {
		data["total_size"] = utils.HumanBytes(models.S.TotalSize)
		data["total_files"] = models.S.TotalFiles
	}
	c.HTML(http.StatusOK, "index.html", data)
}

// SimpleIndex is a simple rendering for webapp purpose
func SimpleIndex(c *gin.Context) {
	logger.InfoC(c, "server", "GET /simple")
	data := gin.H{
		"size_limit":     utils.HumanBytes(uint64(conf.C.SizeLimit * utils.MegaByte)),
		"sensitive_mode": conf.C.SensitiveMode,
	}
	c.HTML(http.StatusOK, "mobile.html", data)
}
