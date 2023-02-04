package restapi

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/dueckminor/mypi-tools/go/util"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	// authURI   string
	webpackDebug  = flag.String("webpack-debug", "", "The debug URI")
	distDir       = flag.String("dist-dir", "", "The location of the dist dir")
	port          = flag.Int("port", 8080, "The port")
	tlsPort       = flag.Int("tlsport", 443, "The TLS port")
	execDebug     = flag.String("exec", "", "start process")
	localhostOnly = flag.Bool("localhost-only", false, "Listen on localhost only")
	listenHost    = ""
)

func prepare(r *gin.Engine) {
	if len(*execDebug) > 0 {
		cmd := exec.Command(*execDebug)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		defer cmd.Wait()
	}

	if len(*webpackDebug) > 0 {
		r.Use(ginutil.SingleHostReverseProxy(*webpackDebug))
	} else {
		usedDistDir := *distDir
		if len(usedDistDir) == 0 {
			usedDistDir = "./dist"
		}
		r.Use(static.ServeRoot("/", usedDistDir))
		r.NoRoute(func(c *gin.Context) {
			c.File(path.Join(usedDistDir, "index.html"))
		})
	}

	if *localhostOnly {
		LocalhostOnly()
	}
}

func LocalhostOnly() {
	listenHost = "localhost"
}

func SetWebpackDebugPort(port int) {
	*webpackDebug = fmt.Sprintf("http://localhost:%d", port)
}

func GetDistDir() (string, bool) {
	if len(*distDir) > 0 {
		if util.FileExists(path.Join(*distDir, "index.html")) {
			return *distDir, true
		}
	}
	return "", false
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
