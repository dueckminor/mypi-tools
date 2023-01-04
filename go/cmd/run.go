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
		if err == nil {
			parsedArgs, err = command.UnmarshalArgs(data)
		}
	} else if len(args) == 2 && args[1][0] == '{' {
		parsedArgs, err = command.UnmarshalArgs([]byte(args[1]))
	} else {
		parsedArgs, err = command.ParseArgs(args[1:])
	}
	if err != nil {
		return true, err
	}
	return true, command.Execute(parsedArgs)
}
