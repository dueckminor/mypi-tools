package setup

import (
	"archive/zip"
	"io"
)

type zipWriter struct {
	zw *zip.Writer
}

func NewZipWriter(w io.Writer) (DirWriter, error) {
	return &zipWriter{zw: zip.NewWriter(w)}, nil
}

func (w *zipWriter) CreateFile(fi FileInfo) (io.WriteCloser, error) {
	fw, err := w.zw.Create(fi.Name)
	if err != nil {
		return nil, err
	}
	return makeWriterNoopCloser(fw), nil
}

func (w *zipWriter) Close() error {
	return w.zw.Close()
}
