package ginutil

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func SingleHostReverseProxy(target string, options ...string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
			r.SetXForwarded()
			r.Out.Host = r.In.Host // if desired
			//fmt.Println(r.Out.Header)
		},
	}

	useExternalHostname := false

	proxy.ModifyResponse = func(resp *http.Response) error {
		location := resp.Header.Get("Location")
		if strings.HasPrefix(location, target) {
			newLocation := location[len(target):]
			resp.Header.Set("Location", newLocation)
		}
		return nil
	}

	for _, option := range options {
		switch option {
		case "insecure":
			proxy.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		case "external-hostname":
			useExternalHostname = true
		}
	}

	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		req := c.Request
		if !useExternalHostname {
			req.Host = url.Hostname()
		}
		proxy.ServeHTTP(c.Writer, req)
	}
}
