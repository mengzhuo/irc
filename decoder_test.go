package irc

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	var err error
	target := []byte(
		`:Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n
:Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n`)
	buf := bytes.NewBuffer(target)
	dec := NewDecoder(buf)
	for i := 0; i < 2; i++ {

		msg := new(Msg)
		err = dec.Decode(msg)

		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(msg.Cmd(), []byte("COMMAND")) {
			t.Error(msg)
		}

	}
}
