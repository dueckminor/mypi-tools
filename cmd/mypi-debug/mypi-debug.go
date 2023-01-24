package main

import (
	"github.com/dueckminor/mypi-tools/go/debug"
	"github.com/dueckminor/mypi-tools/go/ssh"
)

func ConnectToMypiDebug() (err error) {
	client := &ssh.Client{}
	err = client.AddPrivateKeyFile("id_rsa")
	if err != nil {
		return err
	}
	err = client.Dial("pi", "mypi:2022")
	if err != nil {
		return err
	}

	dial := &ssh.DialNet{
		Network: "tcp",
		Address: "127.0.0.1:8443",
	}

	go func() {
		defer client.Close()
		client.RemoteForwardDial("0.0.0.0:8443", dial)
	}()

	return nil
}

func main() {
	err := ConnectToMypiDebug()
	if err != nil {
		panic(err)
	}

	serviceAuth := debug.NewMypiService("mypi-auth")
	serviceAuth.StartGo()

	stop := make(chan bool)
	<-stop
}
