package webhandler

import (
	"github.com/gin-gonic/gin"
)

func getDynDNS(c *gin.Context) {

}

// SetupEnpoints registers the http endpoints
func SetupEndpoints(r *gin.Engine) {
	r.GET("api/dyndns", getDynDNS)
	r.PUT("api/dyndns", putDynDNS)
}
