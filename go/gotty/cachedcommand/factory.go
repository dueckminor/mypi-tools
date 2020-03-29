package cachedcommand

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/dueckminor/mypi-tools/go/gotty/server"
)

type Factory struct {
	name string
}

func NewFactory(name string) (*Factory, error) {
	return &Factory{name: name}, nil
}

func (factory *Factory) Name() string {
	return factory.name
}

type CachedCommand struct {
	name         string
	columns      int
	rows         int
	Pty          *os.File
	Tty          *os.File
	mutex        sync.Mutex
	waitForRead  bool
	waitForReadC chan bool
}

var (
	cachedCommands = make(map[string]*CachedCommand)
)

func New(name string) (c *CachedCommand, err error) {
	if cachedCommand, ok := cachedCommands[name]; ok {
		return cachedCommand, nil
	}
	cachedCommand := &CachedCommand{
		name:         name,
		waitForReadC: make(chan bool),
	}
	cachedCommand.createTty()
	cachedCommands[name] = cachedCommand

	return cachedCommand, nil
}

func (c *CachedCommand) createTty() (err error) {
	if c.Pty == nil {
		c.Pty, c.Tty, err = pty.Open()
		if c.waitForRead {
			c.waitForRead = false
			c.waitForReadC <- true
		}
	}
	return err
}

func AttachProcess(name string, command *exec.Cmd) (err error) {
	cachedCommand, err := New(name)
	if err != nil {
		return err
	}

	cachedCommand.createTty()

	command.Stdout = cachedCommand.Tty
	command.Stderr = cachedCommand.Tty
	command.Env = append(os.Environ(), "TERM=xterm-256color")

	if command.SysProcAttr == nil {
		command.SysProcAttr = &syscall.SysProcAttr{}
	}
	command.SysProcAttr.Setctty = true
	command.SysProcAttr.Setsid = true
	command.SysProcAttr.Ctty = int(cachedCommand.Tty.Fd())

	return nil
}

func (factory *Factory) New(params map[string][]string) (server.Slave, error) {
	return New(factory.name)
}

func (c *CachedCommand) Read(p []byte) (n int, err error) {
	for {
		if c.waitForRead {
			<-c.waitForReadC
			c.waitForRead = false
		}
		n, err = c.Pty.Read(p)
		if err == nil {
			return
		} else if err == io.EOF {
			c.Pty.Close()
			c.Tty.Close()
			c.Pty = nil
			c.Tty = nil
			if n != 0 {
				return n, nil
			}
			c.waitForRead = true
		} else {
			return
		}
	}
}

func (c *CachedCommand) Write(p []byte) (n int, err error) {
	return c.Pty.Write(p)
}

func (c *CachedCommand) Close() error {
	if nil != c.Tty {
		c.Tty.Close()
		c.Tty = nil
	}
	if nil != c.Pty {
		c.Pty.Close()
		c.Pty = nil
	}
	return nil
}

func (c *CachedCommand) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{}
	// return map[string]interface{}{
	// 	"command": lcmd.command,
	// 	"argv":    lcmd.argv,
	// 	"pid":     lcmd.cmd.Process.Pid,
	// }
	return nil
}

// ResizeTerminal sets a new size of the terminal.
func (c *CachedCommand) ResizeTerminal(columns int, rows int) error {
	c.columns = columns
	c.rows = rows

	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(rows),
		uint16(columns),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		c.Pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}
