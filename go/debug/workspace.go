package debug

import (
	"path"
	"runtime"
)

func GetWorkspaceRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "../..")
}
