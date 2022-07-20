package cmdmakesd

import (
	"os"
	"os/exec"

	"github.com/dueckminor/mypi-tools/go/cmd"
)

type cmdStartNode struct{}

type settings struct {
}

func (cmd cmdStartNode) ParseArgs(args []string) (parsedArgs interface{}, err error) {
	settings := settings{}
	return settings, nil
}

func (cmd cmdStartNode) UnmarshalArgs(marshaledArgs []byte) (parsedArgs interface{}, err error) {
	settings := settings{}
	return settings, err
}

func (cmd cmdStartNode) Execute(parsedArgs interface{}) error {
	_, ok := parsedArgs.(settings)
	if !ok {
		return os.ErrInvalid
	}
	cmdStartNode := exec.Command("bash", "-c", "scripts/mypi-auth/start_local_js.sh")
	cmdStartNode.Stdin = os.Stdin
	cmdStartNode.Stdout = os.Stdout
	cmdStartNode.Stderr = os.Stderr
	cmdStartNode.Start()
	cmdStartNode.Wait()
	return nil
}

func init() {
	cmd.Register("setup", cmdStartNode{})
}
