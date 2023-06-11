package https

import (
	"flag"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/gin-gonic/gin"
)

var (
	tlsPort = flag.Int("tlsport", 443, "The TLS port")
)

type HTTPSRunner struct {
	keyFile  string
	certFile string
}

var (
	httpsRunner = &HTTPSRunner{}
)

func SetKeyFiles(keyFile, certFile string) {
	httpsRunner.keyFile = keyFile
	httpsRunner.certFile = certFile
}

func (httpsRunner *HTTPSRunner) Run(r *gin.Engine) error {
	return r.RunTLS(restapi.GetListenHost()+":"+strconv.Itoa(*tlsPort), httpsRunner.certFile, httpsRunner.keyFile)
}

func init() {
	restapi.RegisterRunner(httpsRunner)
}
