package irc

import (
	"bufio"
	"io"
)

type Decoder struct {
	scanner *bufio.Scanner
}

func NewDecoder(r io.Reader) *Decoder {
	scanner := bufio.NewScanner(r)
	return &Decoder{scanner}
}

func (d *Decoder) Decode(msg *Msg) (err error) {

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
