package main

import (
	"fmt"
	"os"
	"syscall"
	"flag"
	"os/signal"
	"github.com/dueckminor/mypi-api/go/config"
	"github.com/dueckminor/mypi-api/go/users"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	mypiRoot = flag.String("mypi-root","","The root of the mypi filesystem")
)

func init() {
	flag.Parse()
	if (mypiRoot != nil && len(*mypiRoot) > 0) {
		config.InitApp(*mypiRoot)
	}
	
	state, _ := terminal.GetState(int(syscall.Stdin))
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		terminal.Restore(int(syscall.Stdin),state)
		fmt.Println("")
		os.Exit(1)
	}()
}

func readPassword(prompt string) string {
	fmt.Print(prompt + ": ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return ""
	}
	return string(bytePassword)
}

func main() {
	username := ""
	password := ""
	if flag.NArg() > 0 {
		username = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		password = flag.Arg(1)
	} else {
		for {
			password = readPassword("Enter Password")
			password2 := readPassword("Enter Again")
			if password == password2 {
				break
			}
		}
	}

	err := users.AddUser(username, password)
	if (err != nil) {
		panic(err)
	}

	fmt.Println(users.CheckPassword(username, password))
}
