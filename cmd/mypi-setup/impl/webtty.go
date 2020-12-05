package impl

import (
	"runtime"

	"github.com/dueckminor/mypi-tools/go/gotty/localcommand"
	"github.com/dueckminor/mypi-tools/go/gotty/server/ginhandler"
	"github.com/gin-gonic/gin"
)

// GetWebTTY return list of hosts
func GetWebTTY(c *gin.Context) {
	sh := "zsh"
	if runtime.GOOS == "linux" && runtime.GOARCH == "arm64" {
		sh = "ash"
	}
	factory, err := localcommand.NewFactory(sh, []string{"-l"}, &localcommand.Options{})
	if err == nil {
		ginhandler.Handler(c, factory)
	}
}
