package http

import (
	"flag"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/gin-gonic/gin"
)

var (
	port = flag.Int("port", 8080, "The port")
)

type HTTPRunner struct {
}

func (httpRunner HTTPRunner) Run(r *gin.Engine) error {
	return r.Run(restapi.GetListenHost() + ":" + strconv.Itoa(*port))
}

func init() {
	restapi.RegisterRunner(HTTPRunner{})
}
