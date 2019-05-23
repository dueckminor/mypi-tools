package ginutil

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func SingleHostReverseProxy(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(resp *http.Response) error {
		location := resp.Header.Get("Location")
		if strings.HasPrefix(location, target) {
			newLocation := location[len(target):]
			resp.Header.Set("Location", newLocation)
		}
		return nil
	}

	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		req := c.Request
		req.Host = url.Hostname()
		proxy.ServeHTTP(c.Writer, req)
	}
}
