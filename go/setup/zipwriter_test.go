package setup

import (
	"io"
	"testing"

	. "github.com/onsi/gomega"
)

func TestZipWriter(t *testing.T) {
	g := NewGomegaWithT(t)

	zw, err := NewZipWriter(io.Discard)
	g.Expect(zw, err).NotTo(BeNil())

	for _, fi := range fileInfoTestTarGz {
		w, err := zw.CreateFile(fi.FileInfo)
		g.Expect(err).To(BeNil())
		if fi.Data != nil {
			g.Expect(w).NotTo(BeNil())
			n, err := w.Write(fi.Data)
			g.Expect(n, err).To(Equal(len(fi.Data)))
		}
	}

	zw.Close()
}
