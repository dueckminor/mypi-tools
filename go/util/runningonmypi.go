package util

var runningOnMypifileExists = FileExists

func IsRunningOnMypi() bool {
	return runningOnMypifileExists("/etc/init.d/mypi-control")
}
