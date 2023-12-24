package util

import (
	"os"
	"strings"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FileIsSafe(filename string) bool {
	if strings.Contains(filename, "..") {
		return false
	}
	if strings.Contains(filename, "/") {
		return false
	}
	return true
}

func FileIsSafePath(filename string) bool {
	return !strings.Contains(filename, "..")
}
