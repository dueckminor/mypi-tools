package restapi

import (
	"flag"
	"path"
	"strings"
	"sync"

	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	// authURI   string
	_dist         = flag.String("dist", "./dist", "The debug URI")
	localhostOnly = flag.Bool("localhost-only", false, "Listen on localhost only")
	listenHost    = ""
)

func prepare(r *gin.Engine) {
	if len(runners) == 0 {
		panic("please import at least one of restapi.http or restapi.https")
	}
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

type Runner interface {
	Run(r *gin.Engine) error
}

var runners = make([]Runner, 0)

func RegisterRunner(runner Runner) {
	runners = append(runners, runner)
}

func GetListenHost() string {
	return listenHost
}

func Run(r *gin.Engine) {
	prepare(r)

	wg := sync.WaitGroup{}
	wg.Add(len(runners))

	for _, runner := range runners {
		go func(runner Runner) {
			err := runner.Run(r)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(runner)
	}

	wg.Wait()
}
