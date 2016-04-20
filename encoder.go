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

func (e *Encoder) AppendByte(b byte) {
	e.buf = append(e.buf, b)
}

func (e *Encoder) Append(p []byte) {
	e.buf = append(e.buf, p...)
}

func (e *Encoder) Encode(msg *Msg) (n int, err error) {
	e.Lock()
	defer e.Unlock()

	e.buf = e.buf[:0]

	if msg.cmd == nil {
		return 0, errors.New("no command")
	}

	if msg.prefix != nil {
		e.AppendByte(prefixSymbol)
		e.Append(msg.prefix)
		e.AppendByte(space)
	}

	e.Append(msg.cmd)

	if msg.paramsCount != 0 {
		e.AppendByte(space)
		for i := 0; i < msg.paramsCount; i++ {
			e.Append(msg.params[i])
			if i != msg.paramsCount-1 {
				e.AppendByte(space)
			}
		}
	}

	if msg.trailing != nil {
		e.AppendByte(space)
		e.AppendByte(prefixSymbol)
		e.Append(msg.trailing)
	}

	e.Append([]byte("\r\n"))
	return e.w.Write(e.buf)
}
