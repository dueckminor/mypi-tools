package cmdmakesd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/fatih/color"
)

// var (
// 	dirSetup = flag.String("dir-setup", "", "the directory containing the setup files")
// 	dirDist  = flag.String("dir-dist", "", "the directory containing the dist files")
// )

type cmdSetup struct{}

type settings struct {
}

func (cmd cmdSetup) ParseArgs(args []string) (parsedArgs interface{}, err error) {
	settings := settings{}
	return settings, nil
}

func (cmd cmdSetup) UnmarshalArgs(marshaledArgs []byte) (parsedArgs interface{}, err error) {
	settings := settings{}
	return settings, err
}

func (cmd cmdSetup) Execute(parsedArgs interface{}) error {
	_, ok := parsedArgs.(settings)
	if !ok {
		return os.ErrInvalid
	}

	c := color.New(color.BgBlue).Add(color.FgHiYellow)

	time.Sleep(time.Second)
	fmt.Println("")
	c.Print("                               ")
	fmt.Println("")
	c.Print(" --- Starting setup-alpine --- ")
	fmt.Println("")
	c.Print("                               ")
	fmt.Println("")
	fmt.Println("")

	cmdSetup := exec.Command("/mypi-control/setup-phase-selector")
	cmdSetup.Stdin = os.Stdin
	cmdSetup.Stdout = os.Stdout
	cmdSetup.Stderr = os.Stderr
	cmdSetup.Start()
	cmdSetup.Wait()
	return nil
}

func init() {
	cmd.Register("setup", cmdSetup{})
}
