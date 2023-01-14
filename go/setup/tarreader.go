package setup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func TarGzFileExtract(tarFile string, w DirWriter) error {
	f, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return TarGzExtract(f, w)
}

func TarGzExtract(r io.Reader, w DirWriter) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	return TarExtract(gzr, w)
}

func TarExtract(r io.Reader, w DirWriter) error {
	tr := NewTarReader(r)

	for {
		fi, err := tr.NextFile()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fw, err := w.CreateFile(fi)
		if err != nil {
			return err
		}
		if fw != nil {
			defer fw.Close()
			fr, err := tr.OpenFile()
			if err != nil {
				return err
			}
			defer fr.Close()
			if _, err = io.Copy(fw, fr); err != nil {
				return err
			}
		}
	}
	return nil
}

type TarReader struct {
	r *tar.Reader
	h *tar.Header
}

func (r *TarReader) NextFile() (fi FileInfo, err error) {
	for {
		r.h, err = r.r.Next()
		if err != nil {
			return
		}

		fi.Name = r.h.Name
		fi.Mode = r.h.Mode
		fi.Size = r.h.Size

		switch r.h.Typeflag {
		case tar.TypeReg:
			fi.Type = FileTypeFile
			return fi, err
		case tar.TypeDir:
			fi.Type = FileTypeDir
			return fi, err
		case tar.TypeSymlink:
			fi.Type = FileTypeSoftlink
			fi.Linkname = r.h.Linkname
			return fi, err
		default:
			r.h = nil
		}
	}
}

func (r *TarReader) OpenFile() (io.ReadCloser, error) {
	if r.h == nil || r.h.Typeflag != tar.TypeReg {
		return nil, fmt.Errorf("TarReader: current file is no regular file")
	}
	return makeReaderNoopCloser(r.r), nil
}

func NewTarReader(r io.Reader) DirSeqReader {
	return &TarReader{r: tar.NewReader(r)}
}
