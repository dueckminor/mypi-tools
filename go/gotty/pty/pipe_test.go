package pty

import (
	"fmt"
	"io"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestPipes(t *testing.T) {
	g := NewGomegaWithT(t)

	w, r := newPipe()

	w.Write([]byte{1, 2, 3})

	b := make([]byte, 2)

	g.Expect(r.Read(b)).To(Equal(2))
	g.Expect(r.Read(b)).To(Equal(1))
}

func TestPipesReadBeforeWrite(t *testing.T) {
	g := NewGomegaWithT(t)

	done := make(chan bool, 1)

	w, r := newPipe()

	go func() {
		b := make([]byte, 2)
		g.Expect(r.Read(b)).To(Equal(2))
		g.Expect(r.Read(b)).To(Equal(1))
		done <- true
	}()

	time.Sleep(time.Millisecond * 50)
	w.Write([]byte{1, 2, 3})

	<-done
}

func testPipesConcurrentReadWrite(t *testing.T, totalSize, writeSize, readSize, writeSleep, readSleep int, waitForEOF bool) {
	g := NewGomegaWithT(t)

	done := make(chan bool, 1)

	w, r := newPipe()

	var readBuffer []byte

	go func() {
		b := make([]byte, readSize)
		for totalRead := 0; waitForEOF || totalRead < totalSize; {
			readNow, err := r.Read(b)
			if err == io.EOF {
				break
			}
			fmt.Println(b[:readNow])
			g.Expect(readNow >= 0, err).To(BeTrue())
			readBuffer = append(readBuffer, b[:readNow]...)
			totalRead += readNow
		}
		done <- true
	}()

	writeBuffer := make([]byte, writeSize)
	for written := 0; written < totalSize; {
		if writeSleep > 0 {
			time.Sleep(time.Duration(writeSleep) * time.Millisecond)
		}

		for i := 0; i < writeSize; i++ {
			writeBuffer[i] = byte(written + i)
		}

		var writtenNow int
		if written+writeSize > totalSize {
			writtenNow, _ = w.Write(writeBuffer[:totalSize-written])
		} else {
			writtenNow, _ = w.Write(writeBuffer)
		}

		written += writtenNow
	}

	w.Close()

	<-done

	g.Expect(readBuffer).To(HaveLen(totalSize))
	for i := 0; i < totalSize; i++ {
		g.Expect(readBuffer[i]).To(Equal(byte(i)))
	}
}

func TestPipesConcurrentReadWrite(t *testing.T) {
	testPipesConcurrentReadWrite(t, 256, 10, 10, 0, 0, false)
	testPipesConcurrentReadWrite(t, 256, 11, 17, 0, 0, false)
	testPipesConcurrentReadWrite(t, 256, 17, 11, 0, 0, false)
}

func TestPipesConcurrentReadWriteWaitForEOF(t *testing.T) {
	testPipesConcurrentReadWrite(t, 256, 11, 17, 0, 0, true)
	testPipesConcurrentReadWrite(t, 256, 17, 11, 0, 0, true)
	testPipesConcurrentReadWrite(t, 256, 10, 10, 0, 0, true)
}
