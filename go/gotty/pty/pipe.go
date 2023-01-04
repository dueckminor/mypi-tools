package pty

import (
	"io"
)

type pipe struct {
	chunks chan [][]byte
	empty  chan bool
}

type pipeWriter struct {
	pipe *pipe
}

type pipeReader struct {
	pipe *pipe
}

func newPipe() (io.WriteCloser, io.ReadCloser) {
	chunks := make(chan [][]byte, 1)
	empty := make(chan bool, 1)
	empty <- true
	pipe := &pipe{chunks: chunks, empty: empty}
	return &pipeWriter{pipe}, &pipeReader{pipe}
}

func (pipe *pipeWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		return 0, nil
	}

	dup := make([]byte, n)
	copy(dup, p)

	var chunks [][]byte
	select {
	case chunks = <-pipe.pipe.chunks:
	case <-pipe.pipe.empty:
	}

	chunks = append(chunks, dup)
	pipe.pipe.chunks <- chunks
	return n, nil
}

func (pipe *pipeWriter) Close() error {
	var chunks [][]byte
	select {
	case chunks = <-pipe.pipe.chunks:
	case <-pipe.pipe.empty:
	}
	chunks = append(chunks, []byte{})
	pipe.pipe.chunks <- chunks
	return nil
}

func (pipe *pipeReader) Read(p []byte) (n int, err error) {
	chunks := <-pipe.pipe.chunks

	want := len(p)
	offset := 0

	if len(chunks) > 0 && len(chunks[0]) == 0 {
		pipe.pipe.chunks <- chunks
		return 0, io.EOF
	}

	for len(chunks) > 0 && want > 0 {
		chunk := chunks[0]
		have := len(chunk)
		if have == 0 {
			break
		}
		if have > want {
			// we have to read only parts of the chunk
			copy(p[offset:offset+want], chunk)
			chunks[0] = chunk[want:]
			want = 0
		} else {
			copy(p[offset:offset+have], chunk)
			chunks = chunks[1:]
			want = want - have
			offset += have
		}
	}

	if len(chunks) == 0 {
		pipe.pipe.empty <- true
	} else {
		pipe.pipe.chunks <- chunks
	}

	return len(p) - want, nil
}

func (pipe *pipeReader) Close() error {
	return nil
}
