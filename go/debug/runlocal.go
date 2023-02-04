package debug

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func runLocal(stdout io.Writer, commandline ...string) (err error) {
	cmd := exec.Command(commandline[0], commandline[1:]...)
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func RunLocal(commandline ...string) (err error) {
	return runLocal(os.Stdout, commandline...)
}

func RunLocalCaptureOutput(commandline ...string) (out string, err error) {
	b := &strings.Builder{}
	err = runLocal(b, commandline...)
	if err != nil {
		return "", err
	}
	return b.String(), err
}
