package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CertInfo struct {
	Text string `json:"text"`
}

// GetCertificates return list of hosts
func GetCertificates(c *gin.Context) {
	c.JSON(http.StatusOK, []CertInfo{CertInfo{"localhost"}, CertInfo{"mypi"}})
}
