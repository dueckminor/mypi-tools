package debug

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"

	"github.com/dueckminor/mypi-tools/go/ssh"
)

type SSHConnector interface {
	io.Closer
	Run(ctx context.Context, uri string, port int, tty io.Writer) (result chan error, err error)
	SetLocalRouterPort(port int) (err error)
	GetFS() (fs.FS, error)
	GetHttpFS() (http.FileSystem, error)
}

type sshConnector struct {
	client *ssh.Client
	dial   *ssh.DialNet
	tty    io.Writer

	httpfs http.FileSystem
}

type httpfs struct {
	conn *sshConnector
}

func (fs httpfs) Open(name string) (f http.File, err error) {
	if fs.conn.httpfs == nil {
		return nil, fmt.Errorf("File not found: %s", name)
	}
	return fs.conn.httpfs.Open(name)
}

func (c *sshConnector) Close() (err error) {
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
	return nil
}

func (c *sshConnector) SetLocalRouterPort(port int) (err error) {
	c.dial.Address = fmt.Sprintf("127.0.0.1:%d", port)
	return nil
}

func (c *sshConnector) Log(a ...any) {
	fmt.Fprint(c.tty, a...)
}
func (c *sshConnector) Logln(a ...any) {
	fmt.Fprintln(c.tty, a...)
}
func (c *sshConnector) Logf(format string, a ...any) {
	fmt.Fprintf(c.tty, format, a...)
}

func (c *sshConnector) Run(ctx context.Context, uri string, port int, tty io.Writer) (result chan error, err error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	username := "pi"
	if parsedURI.User != nil {
		username = parsedURI.User.Username()
	}

	if tty == nil {
		tty = os.Stderr
	}
	c.tty = tty

	c.client = &ssh.Client{}
	err = c.client.AddPrivateKeyFile("id_rsa")
	if err != nil {
		return nil, err
	}
	c.Logf("Username: %v\n", username)
	c.Logf("Host: %v\n", parsedURI.Host)
	err = c.client.Dial(username, parsedURI.Host)
	if err != nil {
		return nil, err
	}

	c.httpfs, err = c.client.GetHttpFS()
	if err != nil {
		return nil, err
	}

	if c.dial == nil {
		c.dial = &ssh.DialNet{
			Network: "tcp",
			Address: fmt.Sprintf("127.0.0.1:%d", port),
		}
	}

	result = make(chan error)

	go func() {
		defer c.client.Close()
		c.Logln("Listen on remote port: 0.0.0.0:8443")
		err := c.client.RemoteForwardDial(ctx, "0.0.0.0:8443", c.dial)
		if err != nil {
			c.Logln("RemoteForwardDial failed:", err)
		}
		result <- err
	}()

	return result, nil
}

func (c *sshConnector) GetFS() (fs.FS, error) {
	return c.client.GetFS()
}
func (c *sshConnector) GetHttpFS() (fs http.FileSystem, err error) {
	return httpfs{c}, nil
}
