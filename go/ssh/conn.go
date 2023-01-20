package ssh

import (
	"io"
	"io/ioutil"
	"net"
	"os/user"
	"path"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	signers []ssh.Signer

	client *ssh.Client
}

func (c *Client) AddPrivateKeyFile(filename string) (err error) {
	if !strings.ContainsAny(filename, "/\\") {
		user, err := user.Current()
		if err != nil {
			return err
		}
		filename = path.Join(user.HomeDir, ".ssh", filename)
	}
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	c.signers = append(c.signers, signer)
	return nil
}

func (c *Client) Dial(username, addr string) (err error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(c.signers...)},
		HostKeyCallback: c.HostKeyCallback,
	}
	c.client, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *Client) Listen(n, addr string) (listener net.Listener, err error) {
	return c.client.Listen(n, addr)
}

func (c *Client) HostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func (c *Client) RemoteForward(remoteAddr, localAddr string) error {
	incoming, err := c.Listen("tcp", remoteAddr)
	if err != nil {
		return err
	}
	defer incoming.Close()

	stop := make(chan bool)

	go func() {
		for {
			remote_conn, err := incoming.Accept()
			if err != nil {
				stop <- true
			}
			go func() {
				defer remote_conn.Close()
				local_conn, err := net.Dial("tcp", localAddr)
				if err != nil {
					return
				}
				defer local_conn.Close()

				done := make(chan bool)

				go func() {
					io.Copy(remote_conn, local_conn)
					<-done
				}()

				io.Copy(local_conn, remote_conn)

				<-done
			}()
		}
	}()

	<-stop

	return nil
}
