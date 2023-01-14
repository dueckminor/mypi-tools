package pty

import (
	"io"
	"os/exec"
)

type Pty interface {
	io.ReadWriteCloser
	SetSize(sx int, sy int) error
	AttachProcess(command *exec.Cmd) error
}

func NewPty() (Pty, error) {
	return newPty()
}
