package setup

import "io"

type DirWriter interface {
	CreateFile(fi FileInfo) (io.WriteCloser, error)
	Close() error
}

//////////////////////////////////////////////////////////////////////////////

type writerNoopCloser struct {
	w io.Writer
}

func (w writerNoopCloser) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}
func (w writerNoopCloser) Close() (err error) {
	return nil
}

func makeWriterNoopCloser(w io.Writer) io.WriteCloser {
	return writerNoopCloser{w}
}

//////////////////////////////////////////////////////////////////////////////
