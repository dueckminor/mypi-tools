package setup

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

type DirReader interface {
	SelectFile(name string) (FileInfo, error)
	OpenFile() (io.ReadCloser, error)
}

type DirSeqReader interface {
	NextFile() (FileInfo, error)
	OpenFile() (io.ReadCloser, error)
}

//////////////////////////////////////////////////////////////////////////////

type readerNoopCloser struct {
	r io.Reader
}

func (r readerNoopCloser) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}
func (w readerNoopCloser) Close() (err error) {
	return nil
}

func makeReaderNoopCloser(r io.Reader) io.ReadCloser {
	return readerNoopCloser{r}
}

//////////////////////////////////////////////////////////////////////////////

func NewStringReader(s string) io.Reader {
	return strings.NewReader(s)
}

//////////////////////////////////////////////////////////////////////////////

// func NewBytesReader(b []byte) io.Reader {
// 	return bytes.NewReader(b)
// }

func NewBase64Reader(r io.Reader) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, r)
}

func NewGZipReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}
