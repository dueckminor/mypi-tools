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

	err = pty.SetSize(20, 10)
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("tput", "lines")
	err = pty.AttachProcess(cmd)
	g.Expect(err).NotTo(HaveOccurred())

	err = cmd.Start()
	g.Expect(err).NotTo(HaveOccurred())

	wg := sync.WaitGroup{}
	wg.Add(1)

	buf := make([]byte, 20)

	go func() {
		defer wg.Done()
		n, err := pty.Read(buf)
		g.Expect(n, err).To(Equal(4))
	}()

	err = cmd.Wait()
	g.Expect(err).NotTo(HaveOccurred())
	wg.Wait()

	g.Expect(buf[:4]).To(Equal([]byte("10\r\n")))

	g.Expect(pty.Close()).ShouldNot(HaveOccurred())
}

func TestPtyCols(t *testing.T) {
	g := NewGomegaWithT(t)

	pty, err := NewPty()
	g.Expect(pty, err).NotTo(BeNil())

	err = pty.SetSize(20, 10)
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("tput", "cols")
	err = pty.AttachProcess(cmd)
	g.Expect(err).NotTo(HaveOccurred())

	err = cmd.Start()
	g.Expect(err).To(BeNil())

	buf := make([]byte, 20)
	n, err := pty.Read(buf)
	g.Expect(n, err).To(Equal(4))

	err = cmd.Wait()
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(buf[:n]).To(Equal([]byte("20\r\n")))

	g.Expect(pty.Close()).ShouldNot(HaveOccurred())
}
