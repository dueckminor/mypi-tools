package pty

import (
	"os/exec"
	"sync"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPtyLines(t *testing.T) {
	g := NewGomegaWithT(t)

	pty, err := NewPty()
	g.Expect(pty, err).NotTo(BeNil())

	pty.SetSize(20, 10)

	cmd := exec.Command("tput", "lines")
	pty.AttachProcess(cmd)

	err = cmd.Start()
	g.Expect(err).To(BeNil())

	wg := sync.WaitGroup{}
	wg.Add(1)

	buf := make([]byte, 20)

	go func() {
		n, err := pty.Read(buf)
		g.Expect(n, err).To(Equal(4))
		wg.Done()
	}()

	cmd.Wait()
	wg.Wait()

	g.Expect(buf[:4]).To(Equal([]byte("10\r\n")))

	g.Expect(pty.Close()).ShouldNot(HaveOccurred())
}

func TestPtyCols(t *testing.T) {
	g := NewGomegaWithT(t)

	pty, err := NewPty()
	g.Expect(pty, err).NotTo(BeNil())

	pty.SetSize(20, 10)

	cmd := exec.Command("tput", "cols")
	pty.AttachProcess(cmd)

	err = cmd.Start()
	g.Expect(err).To(BeNil())

	buf := make([]byte, 20)
	n, err := pty.Read(buf)
	g.Expect(n, err).To(Equal(4))

	cmd.Wait()

	g.Expect(buf[:n]).To(Equal([]byte("20\r\n")))

	g.Expect(pty.Close()).ShouldNot(HaveOccurred())
}
