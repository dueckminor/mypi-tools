package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/dueckminor/mypi-tools/go/mypi/setup"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dueckminor/mypi-tools/go/cmd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdmakesd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdsetup"
)

func main() {
	if len(os.Args) > 1 {
		done, err := cmd.UnmarshalAndExecute(os.Args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if done {
			os.Exit(0)
		}
	}

	flag.Parse()

	r := gin.Default()
	r.Use(cors.Default())

	wh := &webhandler.WebHandler{}

	whDownloads := webhandler.NewWebhandlerDownloads()

	err := wh.SetupEndpoints(r)
	whDownloads.SetupEndpoints(&r.RouterGroup)

	certificates := setup.NewCertificates()
	certificates.CreatePKI()

	webhandler.NewCertHandler(certificates).RegisterEndpoints(r)

	r.GET("/api/hosts", webhandler.GetHosts)

	r.GET("/api/hosts/:host/terminal/webtty", webhandler.MakeForwardToHost(webhandler.GetWebTTY))
	r.GET("/api/hosts/:host/disks", webhandler.MakeForwardToHost(webhandler.GetDisks))

	r.GET("/api/hosts/:host/actions", webhandler.MakeForwardToHost(webhandler.GetActions))
	r.GET("/api/hosts/:host/actions/:action/webtty", webhandler.MakeForwardToHost(webhandler.GetActionWebTTY))

	r.POST("/api/hosts/:host/actions/:action", webhandler.MakeForwardToHost(webhandler.PostAction))

	if err != nil {
		panic(err)
	}

	keyFile := path.Join(os.Getenv("HOME"), ".mypi", "pki", "localhost_tls_priv.pem")
	certFile := path.Join(os.Getenv("HOME"), ".mypi", "pki", "localhost_tls_cert.pem")

	restapi.RunBoth(r, keyFile, certFile)
}
