package debug

import (
	"path"
	"runtime"
)

func GetWorkspaceRoot() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../..")
}
