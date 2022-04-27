package network

import (
	"io"
	"net"
)

func ForwardConn(a, b net.Conn) {
	done := make(chan bool, 2)

	go func() { io.Copy(a, b); done <- true }()
	go func() { io.Copy(b, a); done <- true }()

	<-done
	<-done
}
