package irc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestEncode(t *testing.T) {

	var err error
	var n int
	target := ":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n"

	msg := &Msg{prefix: s2b("Namename!username@hostname"),
		cmd:      s2b("COMMAND"),
		trailing: s2b("Message message message message message"),
	}

	msg.SetParams(bytes.Split(s2b("arg1 arg2 arg3 arg4 arg5 arg6 arg7"), []byte{space})...)

	buf := bytes.NewBuffer([]byte{})
	enc := NewEncoder(buf)
	n, err = enc.Encode(msg)

	if n == 0 || err != nil || buf.String() != target {
		t.Error(n, err, fmt.Sprintf("+%x+", buf.String()))
	}
}

func TestEncodeEmptyMsg(t *testing.T) {
	msg := new(Msg)
	msg.SetParams(bytes.Split(s2b("arg1 arg2 arg3 arg4 arg5 arg6 arg7"), []byte{space})...)
	buf := bytes.NewBuffer([]byte{})
	enc := NewEncoder(buf)
	n, err := enc.Encode(msg)

	if err == nil {
		t.Error(n, err, fmt.Sprintf("+%x+", buf.String()))
	}

}

func BenchmarkEncoder(b *testing.B) {
	msg := new(Msg)
	name := s2b("Namename")
	cmd := s2b("COMMAND")
	trailing := s2b("Message message message message message")
	params := bytes.Split(s2b("arg1 arg2 arg3 arg4 arg5 arg6 arg7"), []byte{space})
	enc := NewEncoder(ioutil.Discard)
	msg.SetParams(params...)
	msg.SetCmd(cmd)
	msg.SetTrailing(trailing)
	msg.SetName(name)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, err := enc.Encode(msg)
		if n == 0 || err != nil {
			b.Fatal(n, err, msg)
		}
	}
}
