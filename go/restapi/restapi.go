package restapi

import (
	"flag"
	"path"
	"strconv"
	"strings"

	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	// authURI   string
	_dist         = flag.String("dist", "./dist", "The debug URI")
	port          = flag.Int("port", 8080, "The port")
	tlsPort       = flag.Int("tlsport", 443, "The TLS port")
	localhostOnly = flag.Bool("localhost-only", false, "Listen on localhost only")
	listenHost    = ""
)

func prepare(r *gin.Engine) {
	if *localhostOnly {
		LocalhostOnly()
	}
	dist := *_dist
	if len(dist) == 0 {
		dist = "./dist"
	}
	if strings.HasPrefix(dist, "http://") || strings.HasPrefix(dist, "https://") {
		r.Use(ginutil.SingleHostReverseProxy(dist))
	} else {
		r.Use(static.ServeRoot("/", dist))
		r.NoRoute(func(c *gin.Context) {
			c.File(path.Join(dist, "index.html"))
		})
	}
}

func LocalhostOnly() {
	listenHost = "localhost"
}

func Run(r *gin.Engine) {
	prepare(r)
	panic(r.Run(listenHost + ":" + strconv.Itoa(*port)))
}

func RunTLS(r *gin.Engine, keyFile, certFile string) {
	prepare(r)
	panic(r.RunTLS(listenHost+":"+strconv.Itoa(*tlsPort), certFile, keyFile))
}

func RunBoth(r *gin.Engine, keyFile, certFile string) {
	prepare(r)

	go func() {
		r.RunTLS(listenHost+":"+strconv.Itoa(*tlsPort), certFile, keyFile)
	}()

	panic(r.Run(listenHost + ":" + strconv.Itoa(*port)))
}
