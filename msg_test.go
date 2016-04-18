package irc

import (
	"bytes"
	"testing"
)

func TestNewMsg(t *testing.T) {
	target := []byte(":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n")
	m, err := NewMsg(target)
	if !bytes.Equal(m.Cmd(), []byte("COMMAND")) || err != nil {
		t.Errorf("failed:%s", m)
	}
}
