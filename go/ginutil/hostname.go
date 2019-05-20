package ginutil

import (
	"github.com/gin-gonic/gin"
)

func GetHostname(c *gin.Context) (hostname string) {
	hostname = c.Request.Header.Get("X-Forwarded-Host")
	if len(hostname) == 0 {
		hostname = c.Request.Host
	}
	return hostname
}

func GetScheme(c *gin.Context) (scheme string) {
	scheme = c.Request.Header.Get("X-Forwarded-Proto")
	if len(scheme) == 0 {
		scheme = "http"
	}
	return scheme
}

func OnHostname(hostname string, f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		if hostname == GetHostname(c) {
			f(c)
		}
	}
}
