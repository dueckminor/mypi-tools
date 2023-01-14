package setup

import (
	"archive/tar"
	"io"
	"os"
	"time"
)

func NewTarWriter(w io.Writer) (tw *TarWriter, err error) {
	return &TarWriter{w: tar.NewWriter(w)}, nil
}

type TarWriter struct {
	w *tar.Writer
}

func (tw *TarWriter) AddBuffer(filename string, data []byte, mode int64) error {
	w, err := tw.CreateFile(filename, mode, int64(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (tw *TarWriter) CreateFile(filename string, mode, size int64) (io.Writer, error) {
	header := new(tar.Header)
	header.Name = filename
	header.Size = size
	header.Mode = mode
	header.ModTime = time.Now()
	// write the header to the tarball archive
	if err := tw.w.WriteHeader(header); err != nil {
		return nil, err
	}
	return tw.w, nil
}

func (tw *TarWriter) AddLink(filename, linkname string, mode int64) error {
	header := new(tar.Header)
	header.Name = filename
	header.Linkname = linkname
	header.Typeflag = tar.TypeSymlink
	header.Mode = int64(mode | int64(os.ModeSymlink))
	header.ModTime = time.Now()
	// write the header to the tarball archive
	return tw.w.WriteHeader(header)
}

func (tw *TarWriter) Close() error {
	return tw.w.Close()
}
