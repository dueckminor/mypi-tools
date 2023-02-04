package debug

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func runOnMypi(stdout io.Writer, commandline ...string) (err error) {
	args := []string{"pi@mypi", "-p", "2022"}
	args = append(args, commandline...)
	cmd := exec.Command("ssh", args...)
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func RunOnMypi(commandline ...string) (err error) {
	return runOnMypi(os.Stdout, commandline...)
}

func RunOnMypiCaptureOutput(commandline ...string) (out string, err error) {
	b := &strings.Builder{}
	err = runOnMypi(b, commandline...)
	if err != nil {
		return "", err
	}
	return b.String(), err
}
