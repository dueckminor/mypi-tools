package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/dueckminor/mypi-api/go/users"
	"golang.org/x/crypto/ssh/terminal"
)

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
	if len(os.Args) > 1 {
		username = os.Args[1]
	}
	if len(os.Args) > 2 {
		password = os.Args[2]
	} else {
		for {
			password = readPassword("Enter Password")
			password2 := readPassword("Enter Again")
			if password == password2 {
				break
			}
		}
	}

	users.AddUser(username, password)

	fmt.Println(users.CheckPassword(username, password))
}
