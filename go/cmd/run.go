package cmd

import (
	"io/ioutil"
	"os"
)

func UnmarshalAndExecute(args []string) (done bool, err error) {
	command, err := GetCommand(args[0])
	if err != nil {
		return false, nil
	}

	var parsedArgs interface{}
	if len(args) == 2 && args[1] == "@" {
		var data []byte
		data, err = ioutil.ReadAll(os.Stdin)
		parsedArgs, err = command.UnmarshalArgs(data)
	} else {
		parsedArgs, err = command.ParseArgs(args[1:])
	}
	if err != nil {
		return true, err
	}
	return true, command.Execute(parsedArgs)
}
