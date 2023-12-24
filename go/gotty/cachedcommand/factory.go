package cachedcommand

import (
	"io"
	"os/exec"

	"github.com/dueckminor/mypi-tools/go/gotty/pty"
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
	Pty          pty.Pty
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
	err = cachedCommand.createTty()
	if err != nil {
		return nil, err
	}
	cachedCommands[name] = cachedCommand

	return cachedCommand, nil
}

func (c *CachedCommand) createTty() (err error) {
	if c.Pty == nil {
		c.Pty, err = pty.NewPty()
		if err != nil {
			return err
		}
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

	err = cachedCommand.createTty()
	if err != nil {
		return err
	}
	return cachedCommand.Pty.AttachProcess(command)
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
			c.Close()
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
}

// ResizeTerminal sets a new size of the terminal.
func (c *CachedCommand) ResizeTerminal(columns int, rows int) error {
	c.columns = columns
	c.rows = rows
	return c.Pty.SetSize(columns, rows)
}
