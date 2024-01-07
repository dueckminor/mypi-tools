package util

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestIsRunningOnMypi(t *testing.T) {
	g := NewWithT(t)

	called := false
	runningOnMypifileExists = func(filename string) bool {
		g.Expect(filename).To(Equal("/etc/init.d/mypi-control"))
		called = true
		return true
	}
	defer func() {
		runningOnMypifileExists = FileExists
	}()

	g.Expect(IsRunningOnMypi()).To(BeTrue())
	g.Expect(called).To(BeTrue())
}
