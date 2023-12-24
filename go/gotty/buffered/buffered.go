package buffered

import (
	"io"
	"os/exec"

	"github.com/dueckminor/mypi-tools/go/gotty/pty"
	"github.com/dueckminor/mypi-tools/go/gotty/server"
)

type BufferedTty interface {
	io.WriteCloser
	GetFactory() server.Factory
	CreatePTY() (pty.Pty, error)
	ClosePTY()
	AttachProcess(command *exec.Cmd) error
}

func NewBufferedTty() (tty BufferedTty, err error) {
	t := &bufferedTty{}

	t.first = make(chan *line, 1)
	t.active = make(chan *line, 1)
	active := &line{}
	t.first <- active
	t.active <- active

	return t, err
}

///////////////////////////////////////////////////////////////////

type bufferedTty struct {
	first     chan *line
	active    chan *line
	activePty pty.Pty
}

type line struct {
	num  int
	data []byte
	next *line
	wait chan struct{}
}

func (t *bufferedTty) GetFactory() server.Factory {
	return &BufferedFactory{t: t}
}

func (t *bufferedTty) AttachProcess(command *exec.Cmd) error {
	pty, err := t.CreatePTY()
	if err != nil {
		return err
	}
	return pty.AttachProcess(command)
}

func (t *bufferedTty) CreatePTY() (pty.Pty, error) {
	if t.activePty != nil {
		return t.activePty, nil
	}

	pty, err := pty.NewPty()
	if err != nil {
		return nil, err
	}

	t.activePty = pty

	go func() {
		defer pty.Close()
		_, err = io.Copy(t, pty)
		if err != nil {
			return
		}
		t.activePty = nil
	}()

	return pty, nil
}

func (t *bufferedTty) ClosePTY() {
	if t.activePty != nil {
		t.activePty.Close()
		t.activePty = nil
	}
}

func (t *bufferedTty) Write(p []byte) (n int, err error) {
	n = len(p)

	iStart := 0
	for i, b := range p {
		if b == '\n' {
			active := <-t.active
			active.data = append(active.data, p[iStart:i+1]...)
			active.next = &line{
				num: active.num + 1,
			}
			if active.wait != nil {
				close(active.wait)
				active.wait = nil
			}

			t.active <- active.next

			iStart = i + 1
		}
	}
	if iStart < len(p) {
		active := <-t.active
		active.data = append(active.data, p[iStart:]...)
		if active.wait != nil {
			close(active.wait)
			active.wait = nil
		}

		t.active <- active
	}
	return n, nil
}

func (t *bufferedTty) Close() (err error) {
	return nil
}

///////////////////////////////////////////////////////////////////////////

type BufferedFactory struct {
	t *bufferedTty
}

func (f *BufferedFactory) New(params map[string][]string) (server.Slave, error) {
	s := &BufferedSlave{
		t:       f.t,
		current: <-f.t.first,
		active:  <-f.t.active,
		pos:     0,
	}
	f.t.first <- s.current
	f.t.active <- s.active
	return s, nil
}

func (f *BufferedFactory) Name() string {
	return "b"
}

///////////////////////////////////////////////////////////////////////////

type BufferedSlave struct {
	t       *bufferedTty
	current *line
	active  *line
	pos     int
}

func (s *BufferedSlave) Close() (err error) {
	s.t = nil
	s.current = nil
	s.active = nil
	s.pos = 0
	return nil
}

func (s *BufferedSlave) Read(p []byte) (n int, err error) {
	m := len(p)
	n = 0
	for {
		isactive := false
		t := s.t
		if t == nil {
			return
		}
		if s.current == s.active {
			// we might read on the active line
			// -> check if it's still active
			// (this locks the active line!)
			s.active = <-t.active
			if s.current == s.active {
				isactive = true
			} else {
				// our current line is no longer the active one
				// -> we can release the lock immediately
				t.active <- s.active
			}
		}

		now := len(s.current.data)
		if now > s.pos {
			now -= s.pos
			if now > m-n {
				now = m - n
			}
			copy(p[n:], s.current.data[s.pos:s.pos+now])
			s.pos += now
			n += now
		}

		if !isactive {
			if n == m {
				return n, nil
			}
			s.current = s.current.next
			s.pos = 0
			continue
		}

		if n > 0 {
			// as we are working on the active line,
			// we finally have to release the lock
			t.active <- s.active
			return n, nil
		}

		if s.active.wait == nil {
			s.active.wait = make(chan struct{})
		}
		wait := s.active.wait
		t.active <- s.active

		<-wait
	}
}

func (s *BufferedSlave) Write(p []byte) (n int, err error) {
	if s.t.activePty != nil {
		return s.t.activePty.Write(p)
	}
	return len(p), nil
}

func (s *BufferedSlave) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{}
}

func (s *BufferedSlave) ResizeTerminal(columns int, rows int) error {
	return nil
}
