package irc

import (
	"bufio"
	"io"
	"sync"
)

type Decoder struct {
	rdr *bufio.Reader
	*sync.Mutex
}

func NewDecoder(r io.Reader) *Decoder {
	rdr := bufio.NewReader(r)
	return &Decoder{rdr, &sync.Mutex{}}
}

// Decode msg from reader
func (d *Decoder) Decode(msg *Msg) (err error) {
	d.Lock()
	defer d.Unlock()

	var line []byte
	line, _, err = d.rdr.ReadLine()
	if err != nil {
		return
	}

	msg.Reset()
	msg.Data = line[:]
	return msg.PeekCmd()
}
