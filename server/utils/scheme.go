package utils

import "github.com/gin-gonic/gin"

// DetectScheme allows to detect the scheme of a request
func DetectScheme(c *gin.Context) string {
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "http" || scheme == "https" {
		return scheme
	}
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}
