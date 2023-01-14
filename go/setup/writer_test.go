package setup

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
)

func TestMakeWriterNoopCloser(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := []byte{}

	bw := bytes.NewBuffer(buf)

	w := makeWriterNoopCloser(bw)

	g.Expect(w.Write([]byte("test"))).To(Equal(4))
	w.Close()

	g.Expect(bw.Bytes()).To(Equal([]byte("test")))
}
