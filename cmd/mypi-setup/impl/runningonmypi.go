package impl

import "github.com/dueckminor/mypi-tools/go/util"

func IsRunningOnMypi() bool {
	return util.FileExists("/etc/init.d/mypi-setup")
}
