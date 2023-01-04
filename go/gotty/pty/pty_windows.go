//go:build windows
// +build windows

package pty

import (
	"io"
	"os/exec"
)

type ptyWindowsToXterm struct {
	out io.Writer
}

func (pty *ptyWindowsToXterm) Write(p []byte) (n int, err error) {
	cNL := 0
	for _, b := range p {
		if b == '\n' {
			cNL += 1
		}
	}
	if cNL == 0 {
		return pty.out.Write(p)
	}

	p2 := make([]byte, len(p)+cNL)
	i := 0
	for _, b := range p {
		if b == '\n' {
			p2[i] = '\r'
			i++
		}
		p2[i] = b
		i++
	}

	n, err = pty.out.Write(p2)
	if n == len(p2) {
		return len(p), err
	} else if n == 0 {
		return 0, err
	} else {
		return len(p), err
	}
}

type ptyWindows struct {
	stdoutWriter io.WriteCloser
	stdoutReader io.ReadCloser
	stdinWriter  io.WriteCloser
	stdinReader  io.ReadCloser
}

func newPty() (Pty, error) {
	stdoutWriter, stdoutReader := newPipe()
	stdinWriter, stdinReader := newPipe()
	return &ptyWindows{
		stdoutWriter: stdoutWriter,
		stdoutReader: stdoutReader,
		stdinWriter:  stdinWriter,
		stdinReader:  stdinReader,
	}, nil
}

func (pty *ptyWindows) Read(p []byte) (n int, err error) {
	return pty.stdoutReader.Read(p)
}

func (pty *ptyWindows) Write(p []byte) (n int, err error) {
	return pty.stdinWriter.Write(p)
}

func (pty *ptyWindows) Close() error {
	pty.stdinWriter.Close()
	pty.stdoutWriter.Close()
	return nil
}

func (pty *ptyWindows) SetSize(sx int, sy int) error {
	return nil
}

func (pty *ptyWindows) AttachProcess(command *exec.Cmd) error {
	command.Stdout = &ptyWindowsToXterm{out: pty.stdoutWriter}
	command.Stderr = &ptyWindowsToXterm{out: pty.stdoutWriter}
	command.Stdin = pty.stdinReader
	// if command.Env == nil {
	// 	command.Env = append([]string{"TERM=xterm-256color"}, os.Environ()...)
	// } else {
	// 	command.Env = append([]string{"TERM=xterm-256color"}, command.Env...)
	// }
	return nil
}
