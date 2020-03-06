package cmd

import "os"

type Cmd interface {
	ParseArgs(args []string) (parsedArgs interface{}, err error)
	UnmarshalArgs(marshaledArgs []byte) (parsedArgs interface{}, err error)
	Execute(parsedArgs interface{}) error
}

var (
	_cmd = make(map[string]Cmd)
)

func Register(name string, cmd Cmd) {
	_cmd[name] = cmd
}

func IsAvailable(name string) bool {
	_, ok := _cmd[name]
	return ok
}

func GetCommand(name string) (command Cmd, err error) {
	cmd, ok := _cmd[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return cmd, nil
}

func ExecuteCmdline(name string, args ...string) error {
	cmd, ok := _cmd[name]
	if !ok {
		return os.ErrNotExist
	}

	parsedArgs, err := cmd.ParseArgs(args)
	if err != nil {
		return err
	}

	return cmd.Execute(parsedArgs)
}

func Execute(name string, args []byte) error {
	cmd, ok := _cmd[name]
	if !ok {
		return os.ErrNotExist
	}

	parsedArgs, err := cmd.UnmarshalArgs(args)
	if err != nil {
		return err
	}

	return cmd.Execute(parsedArgs)
}
