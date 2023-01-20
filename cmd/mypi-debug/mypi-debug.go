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

	go func() {
		defer client.Close()
		client.RemoteForward("0.0.0.0:8443", "127.0.0.1:8443")
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
