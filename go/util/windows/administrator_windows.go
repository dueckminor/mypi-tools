package windows

import "os"

func IsAdminstrator() bool {
	f, _ := os.Open("\\\\.\\PHYSICALDRIVE0")
	if f == nil {
		return false
	}
	f.Close()
	return true
}
