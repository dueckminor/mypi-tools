package ginutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Redirect(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		c.Header("Location", target)
		c.AbortWithStatus(http.StatusFound)
	}
}
