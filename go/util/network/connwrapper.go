package network

import (
	"net"
	"time"
)

type ConnWrapper struct {
	conn      net.Conn
	cacheRead bool
	buff      []byte
}

func NewConnWrapper(conn net.Conn) *ConnWrapper {
	return &ConnWrapper{
		conn:      conn,
		cacheRead: true,
	}
}

func (w *ConnWrapper) Detach() {
	w.conn = nil
}

func (w *ConnWrapper) StopCaching() {
	w.cacheRead = false
}

func (w *ConnWrapper) GetCacheBuffer() []byte {
	return w.buff
}

func (w *ConnWrapper) Read(b []byte) (n int, err error) {
	if w.conn == nil {
		return 0, nil
	}
	n, err = w.conn.Read(b)
	if w.cacheRead && n > 0 {
		w.buff = append(w.buff, b[0:n]...)
	}
	return
}
func (w *ConnWrapper) Write(b []byte) (n int, err error) {
	if w.conn == nil {
		return len(b), nil
	}
	return w.conn.Write(b)
}
func (w *ConnWrapper) Close() error {
	if w.conn == nil {
		return nil
	}
	return w.conn.Close()
}
func (w *ConnWrapper) LocalAddr() net.Addr {
	if w.conn == nil {
		return nil
	}
	return w.conn.LocalAddr()
}
func (w *ConnWrapper) RemoteAddr() net.Addr {
	if w.conn == nil {
		return nil
	}
	return w.conn.RemoteAddr()
}
func (w *ConnWrapper) SetDeadline(t time.Time) error {
	if w.conn == nil {
		return nil
	}
	return w.conn.SetDeadline(t)
}
func (w *ConnWrapper) SetReadDeadline(t time.Time) error {
	if w.conn == nil {
		return nil
	}
	return w.conn.SetReadDeadline(t)
}
func (w *ConnWrapper) SetWriteDeadline(t time.Time) error {
	if w.conn == nil {
		return nil
	}
	return w.conn.SetWriteDeadline(t)
}
