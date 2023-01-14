package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fileWriter struct {
	baseDir string
}

func (w *fileWriter) CreateFile(fi FileInfo) (io.WriteCloser, error) {
	if strings.Contains(fi.Name, "..") {
		return nil, fmt.Errorf("filenames with .. are not allowed")
	}
	absFilename := filepath.Join(w.baseDir, fi.Name)
	absDirname := filepath.Dir(absFilename)
	err := os.MkdirAll(absDirname, os.ModePerm)
	if err != nil {
		return nil, err
	}
	switch fi.Type {
	case FileTypeFile:
		return os.Create(absFilename)
	case FileTypeDir:
		return nil, os.MkdirAll(absFilename, os.ModePerm)
	}
	return nil, nil
}

func (w *fileWriter) Close() error {
	return nil
}

func NewFileWriter(dir string) (DirWriter, error) {
	return &fileWriter{
		baseDir: dir,
	}, nil
}
