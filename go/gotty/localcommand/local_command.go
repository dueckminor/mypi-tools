package localcommand

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/dueckminor/mypi-tools/go/gotty/pty"
	"github.com/pkg/errors"
)

const (
	DefaultCloseSignal  = syscall.SIGINT
	DefaultCloseTimeout = 10 * time.Second
)

type LocalCommand struct {
	command string
	argv    []string

	closeSignal  syscall.Signal
	closeTimeout time.Duration

	cmd       *exec.Cmd
	pty       pty.Pty
	ptyClosed chan struct{}
}

func New(command string, argv []string, options ...Option) (*LocalCommand, error) {
	cmd := exec.Command(command, argv...)

	pty, err := pty.NewPty()
	if err != nil {
		// todo close cmd?
		return nil, errors.Wrapf(err, "failed to start command `%s`", command)
	}

	err = pty.AttachProcess(cmd)
	if err != nil {
		// todo close cmd?
		return nil, errors.Wrapf(err, "failed to start command `%s`", command)
	}

	cmd.Env = []string{"TERM=xterm-256color"}
	cmd.Env = append(cmd.Env, os.Environ()...)

	err = cmd.Start()
	if err != nil {
		// todo close cmd?
		return nil, errors.Wrapf(err, "failed to start command `%s`", command)
	}
	ptyClosed := make(chan struct{})

	lcmd := &LocalCommand{
		command: command,
		argv:    argv,

		closeSignal:  DefaultCloseSignal,
		closeTimeout: DefaultCloseTimeout,

		cmd:       cmd,
		pty:       pty,
		ptyClosed: ptyClosed,
	}

	for _, option := range options {
		option(lcmd)
	}

	// When the process is closed by the user,
	// close pty so that Read() on the pty breaks with an EOF.
	go func() {
		defer func() {
			lcmd.pty.Close()
			close(lcmd.ptyClosed)
		}()

		lcmd.cmd.Wait() // nolint: errcheck
	}()

	return lcmd, nil
}

func (lcmd *LocalCommand) Read(p []byte) (n int, err error) {
	return lcmd.pty.Read(p)
}

func (lcmd *LocalCommand) Write(p []byte) (n int, err error) {
	return lcmd.pty.Write(p)
}

func (lcmd *LocalCommand) Close() error {
	if lcmd.cmd != nil && lcmd.cmd.Process != nil {
		err := lcmd.cmd.Process.Signal(lcmd.closeSignal)
		if err != nil {
			return err
		}
	}
	for {
		select {
		case <-lcmd.ptyClosed:
			return nil
		case <-lcmd.closeTimeoutC():
			err := lcmd.cmd.Process.Signal(syscall.SIGKILL)
			if err != nil {
				return err
			}
		}
	}
}

func (lcmd *LocalCommand) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{
		"command": lcmd.command,
		"argv":    lcmd.argv,
		"pid":     lcmd.cmd.Process.Pid,
	}
}

func (lcmd *LocalCommand) ResizeTerminal(width int, height int) error {
	return lcmd.pty.SetSize(width, height)
}

func (lcmd *LocalCommand) closeTimeoutC() <-chan time.Time {
	if lcmd.closeTimeout >= 0 {
		return time.After(lcmd.closeTimeout)
	}

	return make(chan time.Time)
}
