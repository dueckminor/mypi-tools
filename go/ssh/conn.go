package ssh

import (
	"context"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/pkg/sftp"
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
	key, err := os.ReadFile(filename)
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

type Dial interface {
	Dial() (io.ReadWriteCloser, error)
}

type DialNet struct {
	Network string
	Address string
}

func (d DialNet) Dial() (io.ReadWriteCloser, error) {
	network, address := d.Network, d.Address
	if len(address) == 0 {
		return nil, net.ErrClosed
	}
	return net.Dial(network, address)
}

func (c *Client) RemoteForward(ctx context.Context, remoteAddr, localAddr string) error {
	return c.RemoteForwardDial(ctx, remoteAddr, &DialNet{"tcp", localAddr})
}

func handleRemoteConn(remote_conn net.Conn, dial Dial) {
	defer remote_conn.Close()
	local_conn, err := dial.Dial()
	if err != nil {
		return
	}
	defer local_conn.Close()

	done := make(chan bool)

	go func() {
		io.Copy(remote_conn, local_conn) // nolint: errcheck
		done <- true
	}()

	io.Copy(local_conn, remote_conn) // nolint: errcheck

	<-done
}

func (c *Client) RemoteForwardDial(ctx context.Context, remoteAddr string, dial Dial) error {
	incoming, err := c.Listen("tcp", remoteAddr)
	if err != nil {
		return err
	}
	defer incoming.Close()

	stop := make(chan bool)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			remote_conn, err := incoming.Accept()
			if err != nil {
				break
			}
			go handleRemoteConn(remote_conn, dial)
		}
		cancel()
		stop <- true
	}()

	<-ctx.Done()
	incoming.Close()
	cancel()

	<-stop

	return nil
}

func (c *Client) GetFS() (fs fs.FS, err error) {
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return nil, err
	}
	return &FS{sftpClient: sftpClient}, nil
}

func (c *Client) GetHttpFS() (fs http.FileSystem, err error) {
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return nil, err
	}
	return &HttpFS{sftpClient: sftpClient}, nil
}

type HttpFS struct {
	sftpClient *sftp.Client
}

type FS struct {
	sftpClient *sftp.Client
}

func open(sftpClient *sftp.Client, name string) (*FSFile, error) {
	stat, err := sftpClient.Stat(name)
	if err != nil {
		return nil, err
	}
	fsf := &FSFile{stat: stat, name: name, sftpClient: sftpClient}
	if !stat.IsDir() {
		fsf.file, err = sftpClient.Open(fsf.name)
		if err != nil {
			return nil, err
		}
		return fsf, nil
	}

	fsf.entries, err = sftpClient.ReadDir(name)
	if err != nil {
		return nil, err
	}

	return fsf, nil
}

func (fs *HttpFS) Open(name string) (http.File, error) {
	return open(fs.sftpClient, name)
}

func (fs *FS) Open(name string) (fs.File, error) {
	return open(fs.sftpClient, name)
}

type FSFile struct {
	stat       os.FileInfo
	name       string
	sftpClient *sftp.Client
	file       *sftp.File
	entries    []fs.FileInfo
}

func (fsf *FSFile) Stat() (fs.FileInfo, error) {
	return fsf.stat, nil
}
func (fsf *FSFile) Read(b []byte) (int, error) {
	return fsf.file.Read(b)
}
func (fsf *FSFile) Seek(offset int64, whence int) (int64, error) {
	return fsf.file.Seek(offset, whence)
}

func (fsf *FSFile) Close() error {
	if nil != fsf.file {
		err := fsf.file.Close()
		fsf.file = nil
		return err
	}
	return nil
}

type FSDirEntry struct {
	fileinfo fs.FileInfo
}

func (fsde FSDirEntry) Name() string {
	return fsde.fileinfo.Name()
}
func (fsde FSDirEntry) IsDir() bool {
	return fsde.fileinfo.IsDir()
}
func (fsde FSDirEntry) Type() fs.FileMode {
	return fsde.fileinfo.Mode()
}
func (fsde FSDirEntry) Info() (fs.FileInfo, error) {
	return fsde.fileinfo, nil
}

func (fsf *FSFile) Readdir(n int) ([]fs.FileInfo, error) {
	if nil == fsf.entries {
		return nil, io.EOF
	}
	if n <= 0 || len(fsf.entries) <= n {
		result := fsf.entries
		fsf.entries = nil
		return result, nil
	}

	result := fsf.entries[:n]
	fsf.entries = fsf.entries[n:]
	return result, nil
}

func (fsd *FSFile) ReadDir(n int) ([]fs.DirEntry, error) {
	entries, err := fsd.Readdir(n)
	if err != nil {
		return nil, err
	}
	result := make([]fs.DirEntry, len(entries))
	for i, entry := range entries {
		fsde := FSDirEntry{fileinfo: entry}
		result[i] = fsde
	}
	return result, nil
}
