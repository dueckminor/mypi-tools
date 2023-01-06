//go:build !windows
// +build !windows

package pty

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	unix_pty "github.com/creack/pty"
)

type ptyUnix struct {
	Pty *os.File
	Tty *os.File
}

func newPty() (result Pty, err error) {
	pty := &ptyUnix{}
	pty.Pty, pty.Tty, err = unix_pty.Open()
	if err != nil {
		return nil, err
	}
	return pty, nil
}

func (pty *ptyUnix) Read(p []byte) (n int, err error) {
	return pty.Pty.Read(p)
}

func (pty *ptyUnix) Write(p []byte) (n int, err error) {
	return pty.Pty.Write(p)
}

func (pty *ptyUnix) Close() error {
	pty.Pty.Close()
	pty.Pty = nil
	pty.Tty.Close()
	pty.Tty = nil
	return nil
}

func (pty *ptyUnix) SetSize(sx int, sy int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(sy),
		uint16(sx),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		pty.Tty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return fmt.Errorf("Syscall TIOCSWINSZ failed with errno %d", errno)
	}
	return nil
}

func (pty *ptyUnix) AttachProcess(command *exec.Cmd) error {

	command.Stdout = pty.Tty
	command.Stderr = pty.Tty
	command.Stdin = pty.Tty
	command.Env = append(os.Environ(), "TERM=xterm-256color")

	if command.SysProcAttr == nil {
		command.SysProcAttr = &syscall.SysProcAttr{}
	}
	command.SysProcAttr.Setctty = true
	command.SysProcAttr.Setsid = true
	command.SysProcAttr.Ctty = 3
	command.ExtraFiles = []*os.File{pty.Tty}

	return nil
}
