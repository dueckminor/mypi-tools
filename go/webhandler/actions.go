package webhandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/dueckminor/mypi-tools/go/gotty/cachedcommand"
	"github.com/dueckminor/mypi-tools/go/gotty/server/ginhandler"
	"github.com/gin-gonic/gin"
)

func GetActions(c *gin.Context) {
	actions, err := cmd.GetCommands()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, actions)
}

func GetActionWebTTY(c *gin.Context) {
	action := c.Params.ByName("action")
	if _, err := cmd.GetCommand(action); err == nil {
		factory, err := cachedcommand.NewFactory(action)
		if err == nil {
			ginhandler.Handler(c, factory)
		}
	} else {
		fmt.Fprintln(os.Stderr, "action '"+action+"' not found")
	}
}

func GetAction(c *gin.Context) {
	action := c.Params.ByName("action")
	if command, err := cmd.GetCommand(action); err == nil {
		data, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		parsedArgs, err := command.UnmarshalArgs(data)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		data, err = json.Marshal(parsedArgs)

		c := exec.Command(os.Args[0], action, "@")
		err = cachedcommand.AttachProcess(action, c)
		c.Stdin = bytes.NewReader(data)

		go func() {
			c.Start()
			c.Wait()
		}()
	}
}
