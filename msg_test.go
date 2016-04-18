package irc

import (
	"fmt"
	"testing"
)

func TestMsgParse(t *testing.T) {
	m := new(Msg)
	target := []byte(":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n")
	m.Data = target
	m.parseAll()
	//t.Errorf("failed:%s", m)
	fmt.Println(m)
}

var cmd []byte

func BenchmarkParseCmd_long(b *testing.B) {

	target := []byte(":syrk!kalt@millennium.stealth.net QUIT :Gone to have lunch")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMsg(target)
		m.parseCmd()
	}
}

var params [][]byte

func BenchmarkParseAllLong(b *testing.B) {

	target := []byte(":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n")
	b.ResetTimer()
	m := NewMsg(target)
	for i := 0; i < b.N; i++ {
		m.parseAll()
		params = m.Params
	}

}
