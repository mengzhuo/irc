package irc

import (
	"errors"
	"io"
	"sync"
)

var DefaultEncoderBufferrSize = 1024

type Encoder struct {
	w   io.Writer
	buf []byte
	*sync.Mutex
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w,
		make([]byte, DefaultEncoderBufferrSize),
		&sync.Mutex{}}
}

func (e *Encoder) appendByte(b byte) {
	e.buf = append(e.buf, b)
}

func (e *Encoder) append(p []byte) {
	e.buf = append(e.buf, p...)
}

// Encode msg into writer
func (e *Encoder) Encode(msg *Msg) (n int, err error) {
	e.Lock()
	defer e.Unlock()

	e.buf = e.buf[:0]

	if msg.cmd == nil {
		return 0, errors.New("no command")
	}

	if msg.prefix != nil {
		e.appendByte(prefixSymbol)
		e.append(msg.prefix)
		e.appendByte(space)
	}

	e.append(msg.cmd)

	if msg.paramsCount != 0 {
		e.appendByte(space)
		for i := 0; i < msg.paramsCount; i++ {
			e.append(msg.params[i])
			if i != msg.paramsCount-1 {
				e.appendByte(space)
			}
		}
	}

	if msg.trailing != nil {
		e.appendByte(space)
		e.appendByte(prefixSymbol)
		e.append(msg.trailing)
	}

	e.append([]byte("\r\n"))
	return e.w.Write(e.buf)
}
