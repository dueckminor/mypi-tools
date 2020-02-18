package cmd

import "os"

type Cmd interface {
	Execute(args []string) error
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

func Execute(name string, args ...string) error {
	cmd, ok := _cmd[name]
	if !ok {
		return os.ErrNotExist
	}
	return cmd.Execute(args)
}
