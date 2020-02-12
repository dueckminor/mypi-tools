package sudo

import (
	"os"
	"os/exec"
	"strings"
)

// Exec executes a process
func Exec(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

// Sudo checks if it's possible to run a process as root without prompt
func Sudo(password string) error {
	c := exec.Command("sudo", "-S", truePath)
	c.Stdin = strings.NewReader(password + "\n")
	//c.Stdout = os.Stdout
	//c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
