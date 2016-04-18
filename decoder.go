package irc

import (
	"bufio"
	"io"
)

type Decoder struct {
	r       io.Reader
	scanner *bufio.Scanner
}

func NewDecoder(r io.Reader) *Decoder {
	scanner := bufio.NewScanner(r)
	return &Decoder{r, scanner}
}

func (d *Decoder) Decode(msg *Msg) error {

	for d.scanner.Scan() {
		msg.Reset()
		copy(d.scanner.Bytes(), msg.Data)
		msg.parseCmd()
	}

	if d.scanner.Err() != nil {
		return d.scanner.Err()
	}

	return nil
}
