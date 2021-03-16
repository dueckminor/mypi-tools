package util

func IsRunningOnMypi() bool {
	return FileExists("/etc/init.d/mypi-control")
}
