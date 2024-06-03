package webhandler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ForwardToHost(c *gin.Context, f func(c *gin.Context)) {
	host := c.Params.ByName("host")
	req := c.Request

	local := (host == "localhost") || (host == "127.0.0.1")
	if !local {
		if osHost, err := os.Hostname(); err == nil {
			local = osHost == host
		}
	}
	if !local {
		urlHostParts := strings.Split(req.Host, ":")
		if len(urlHostParts) > 0 {
			local = urlHostParts[0] == host
		}
	}

	if local {
		f(c)
	} else {
		newHost := host + ":8080"
		targetURI, _ := url.ParseRequestURI("http://" + newHost)
		proxy := httputil.NewSingleHostReverseProxy(targetURI)
		proxy.ServeHTTP(c.Writer, req)
	}
}

func MakeForwardToHost(f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ForwardToHost(c, f)
	}
}

// GetHosts return list of hosts
func GetHosts(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"localhost", "mypi", "rpi5"})
}
