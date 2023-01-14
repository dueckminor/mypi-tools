package pty

import (
	"os/exec"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPty(t *testing.T) {
	g := NewGomegaWithT(t)

	pty, err := newPty()
	g.Expect(pty, err).NotTo(BeNil())

	pty.SetSize(20, 10)

	cmd := exec.Command("tput", "lines")
	pty.AttachProcess(cmd)

	err = cmd.Run()
	g.Expect(err).To(BeNil())

	buf := make([]byte, 20)
	n, err := pty.Read(buf)
	g.Expect(n, err).To(Equal(4))

	g.Expect(buf[:n]).To(Equal([]byte("10\r\n")))

	cmd = exec.Command("tput", "cols")
	pty.AttachProcess(cmd)

	err = cmd.Run()
	g.Expect(err).To(BeNil())

	n, err = pty.Read(buf)
	g.Expect(n, err).To(Equal(4))

	g.Expect(buf[:n]).To(Equal([]byte("20\r\n")))

	cmd = exec.Command("cat")
	pty.AttachProcess(cmd)

	go func() {
		pty.Write([]byte("test\r\n\x04"))
	}()

	err = cmd.Run()
	g.Expect(err).To(BeNil())

	n, err = pty.Read(buf)
	g.Expect(n, err).To(Equal(16))

	g.Expect(buf[:n]).To(Equal([]byte("test\r\n\r\ntest\r\n\r\n")))

	err = pty.Close()
	g.Expect(err).To(BeNil())
}
