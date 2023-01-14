package setup

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
)

func TestMakeReaderNoopCloser(t *testing.T) {
	g := NewGomegaWithT(t)

	br := bytes.NewBuffer([]byte("test"))

	r := makeReaderNoopCloser(br)

	buf := make([]byte, 4)
	g.Expect(r.Read(buf)).To(Equal(4))
	g.Expect(r.Close()).To(BeNil())

	g.Expect(buf).To(Equal([]byte("test")))
}
