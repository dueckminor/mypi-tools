package webhandler

import (
	"net/http"

	"github.com/dueckminor/mypi-tools/go/services"

	"github.com/gin-gonic/gin"
)

func (wh *WebHandler) getServices(c *gin.Context) {
	result, err := services.GetServices(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, result)
}
