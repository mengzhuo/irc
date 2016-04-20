package irc

import (
	"bufio"
	"io"
	"sync"
)

type Decoder struct {
	scanner *bufio.Scanner
	*sync.Mutex
}

func NewDecoder(r io.Reader) *Decoder {
	scanner := bufio.NewScanner(r)
	return &Decoder{scanner, &sync.Mutex{}}
}

// Decode msg from reader
func (d *Decoder) Decode(msg *Msg) (err error) {
	d.Lock()
	defer d.Unlock()

	for d.scanner.Scan() {
		msg.Reset()
		msg.Data = d.scanner.Bytes()[:]
		err = msg.PeekCmd()
		if err != nil {
			return err
		}
		break
	}

	if d.scanner.Err() != nil {
		return d.scanner.Err()
	}

	return nil
}
