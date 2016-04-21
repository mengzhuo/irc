package irc

import (
	"bytes"
	"fmt"
	"testing"
)

var target = []byte(
	`:Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n
: \r\n`)

func TestDecode(t *testing.T) {
	var err error
	buf := bytes.NewBuffer(target)
	dec := NewDecoder(buf)
	for i := 0; i < 3; i++ {

		msg := new(Msg)
		err = dec.Decode(msg)

		if i == 0 && err != nil && !bytes.Equal(msg.Cmd(), []byte("COMMAND")) {
			t.Error(msg)
		}
		if i == 1 && err == nil {
			t.Error(msg)
		}
		if i == 2 && (msg.cmd != nil || err != nil) {
			// EOF
			t.Error(msg, err)
		}
	}
}

func ExampleDecoder() {
	buf := bytes.NewBuffer([]byte(`:hello NICK world :irc`))
	dec := NewDecoder(buf)
	msg := new(Msg)
	_ = dec.Decode(msg)
	fmt.Println(msg.Cmd()) // NICK
}

func BenchmarkDecoder(b *testing.B) {
	target := []byte(
		`:Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n`)
	buf := bytes.NewBuffer(target)
	dec := NewDecoder(buf)
	msg := new(Msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := dec.Decode(msg)
		if err != nil || msg.Cmd() == nil {
			b.Fatal(err, msg, msg.Data)
		}
	}
}
