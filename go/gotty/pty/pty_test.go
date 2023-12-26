package pty

import (
	"io"
	"os/exec"
	"sync"
	"testing"

	. "github.com/onsi/gomega"
)

func readSomeData(reader io.Reader, buf []byte, atLeast int) (read int, err error) {
	read = 0
	for read < atLeast {
		var n int
		n, err = reader.Read(buf[read:])
		read += n
		if err != nil {
			break
		}
	}
	return read, err
}

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
		read, err := readSomeData(pty, buf, 4)
		g.Expect(read, err).To(Equal(4))
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
	n, err := readSomeData(pty, buf, 4)
	g.Expect(n, err).To(Equal(4))

	err = cmd.Wait()
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(buf[:n]).To(Equal([]byte("20\r\n")))

	g.Expect(pty.Close()).ShouldNot(HaveOccurred())
}
