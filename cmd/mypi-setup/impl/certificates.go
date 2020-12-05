package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCertificates return list of hosts
func GetCertificates(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"localhost", "mypi"})
}
