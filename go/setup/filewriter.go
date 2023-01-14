package setup

import (
	"io"
	"os"
	"path/filepath"
)

type fileWriter struct {
	mountPoint string
}

func (w *fileWriter) CreateFile(fi FileInfo) (io.WriteCloser, error) {
	absFilename := filepath.Join(w.mountPoint, fi.Name)
	absDirname := filepath.Dir(absFilename)
	err := os.MkdirAll(absDirname, os.ModePerm)
	if err != nil {
		return nil, err
	}
	if fi.Type == FileTypeFile {
		return os.Create(absFilename)
	}
	return nil, nil
}

func (w *fileWriter) Close() error {
	return nil
}

func NewFileWriter(dir string) (DirWriter, error) {
	return &fileWriter{
		mountPoint: dir,
	}, nil
}
