package ginhandler

import (
	"github.com/dueckminor/mypi-tools/go/gotty/server"
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context, factory server.Factory) {
	server.Handler(c, c.Writer, c.Request, factory)
}
