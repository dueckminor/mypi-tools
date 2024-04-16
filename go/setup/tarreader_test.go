package setup

import (
	"io"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	testTarGz            = "H4sIAAAAAAAAA+3WTQ6CMBAF4Fl7ip4A+jed87RbjSRYE49vRVYs0E0Lyvs2JUDCNC99IfZUnS6EeVqL5TpdG++CZsuBHWlT7jAprj8a0f2W46gUjcOQ19779PxHxT52+VF3Z6+Ag/cr+dtl/sKBlK461ezo+Z+2ngC2lPbW/2Huf0H/t5D6tK/+f+cvwaD/W0jo/0OL3eV6rvyNqf9Fvuh/Z43Ycv6NaEfKtvg5Pfj5BwAAAAAAAAAAAAAAgP/yBNkg120AKAAA"
	testTarGzWithDevNull = "H4sIAAAAAAAAA+3PMQ7CMBAEQD/FP8DGxnkPEnQRkULC+4kCNBRQhYaZZk+6K/ZO59vuMvd92FBatNbWXLznOudaWjqkUrscUt7XroRYtiz1Ml+n4xhjGIdh+nT3bf/4Iz/zJ9UBAAAAAAAAAAD4X3c6cnQAACgAAA=="
)

var (
	fileInfoTestTarGz = []fileInfoWithData{
		{
			FileInfo: FileInfo{
				Type: FileTypeDir,
				Name: "a/",
				Mode: 0755,
			},
		},
		{
			FileInfo: FileInfo{
				Type: FileTypeFile,
				Name: "a/a.txt",
				Size: 2,
				Mode: 0644,
			},
			Data: []byte("a\n"),
		},
		{
			FileInfo: FileInfo{
				Type: FileTypeDir,
				Name: "b/",
				Mode: 0755,
			},
		},
		{
			FileInfo: FileInfo{
				Type: FileTypeFile,
				Name: "b/b.txt",
				Size: 2,
				Mode: 0644,
			},
			Data: []byte("b\n"),
		},
		{
			FileInfo: FileInfo{
				Type:     FileTypeSoftlink,
				Name:     "a.lnk",
				Linkname: "a/a.txt",
				Mode:     0777,
			},
		},
	}
)

func TestTarReader(t *testing.T) {
	g := NewGomegaWithT(t)

	r, err := NewGZipReader(NewBase64Reader(NewStringReader(testTarGz)))
	g.Expect(r, err).NotTo(BeNil())

	tr := NewTarReader(r)
	g.Expect(tr).NotTo(BeNil())

	fi, err := tr.NextFile()
	g.Expect(err).To(BeNil())
	g.Expect(fi.Type).To(Equal(FileTypeDir))
	g.Expect(fi.Name).To(Equal("a/"))
	g.Expect(fi.Mode).To(Equal(int64(0755)))

	g.Expect(tr.OpenFile()).Error().ShouldNot(BeNil())

	fi, err = tr.NextFile()
	g.Expect(err).To(BeNil())
	g.Expect(fi.Type).To(Equal(FileTypeFile))
	g.Expect(fi.Name).To(Equal("a/a.txt"))
	g.Expect(fi.Mode).To(Equal(int64(0644)))
	g.Expect(fi.Size).To(Equal(int64(2)))

	fi, err = tr.NextFile()
	g.Expect(err).To(BeNil())
	g.Expect(fi.Type).To(Equal(FileTypeDir))
	g.Expect(fi.Name).To(Equal("b/"))
	g.Expect(fi.Mode).To(Equal(int64(0755)))

	fi, err = tr.NextFile()
	g.Expect(err).To(BeNil())
	g.Expect(fi.Type).To(Equal(FileTypeFile))
	g.Expect(fi.Name).To(Equal("b/b.txt"))
	g.Expect(fi.Mode).To(Equal(int64(0644)))
	g.Expect(fi.Size).To(Equal(int64(2)))

	f, err := tr.OpenFile()
	g.Expect(f, err).NotTo(BeNil())

	b, err := io.ReadAll(f)
	g.Expect(b, err).To(Equal([]byte("b\n")))

	fi, err = tr.NextFile()
	g.Expect(err).To(BeNil())
	g.Expect(fi.Type).To(Equal(FileTypeSoftlink))
	g.Expect(fi.Name).To(Equal("a.lnk"))
	g.Expect(fi.Linkname).To(Equal("a/a.txt"))
	g.Expect(fi.Mode).To(Equal(int64(0777)))

	_, err = tr.NextFile()
	g.Expect(err).To(Equal(io.EOF))
}
