package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/conf"
)

// ServerName returns the configured hostname (and port) of the server
func ServerName() string {
	ns := conf.C.NameServer
	if conf.C.AppendPort {
		ns = fmt.Sprintf("%s:%d", conf.C.NameServer, conf.C.Port)
	}
	return ns
}

// ServerURI returns the full URI to the server including scheme
func ServerURI(c *gin.Context) string {
	return DetectScheme(c) + "://" + ServerName()
}
