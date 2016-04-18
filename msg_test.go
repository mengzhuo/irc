package irc

import (
	"bytes"
	"testing"
)

func s2b(s string) []byte {
	return []byte(s)
}

var messageTests = [...]*struct {
	rawMsg   string
	parsed   *Msg
	hostmask bool
	server   bool
}{
	{
		":syrk!kalt@millennium.stealth.net QUIT :Gone to have lunch",
		&Msg{cmd: s2b("QUIT"),
			trailing: s2b("Gone to have lunch"),
			name:     s2b("syrk"),
			user:     s2b("kalt"),
			host:     s2b("millennium.stealth.net"),
		},
		true,
		false,
	},
}

func TestMsgCmd(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.cmd, p.cmd) {
			t.Errorf("failed:%s parsed:%s", z, m.String())
		}
	}
}

func TestMsgPrefixName(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.name, p.name) {
			t.Errorf("failed:%s parsed:%s", z, m.String())
		}
	}
}

func TestMsgUser(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.user, p.user) {
			t.Errorf("failed:%s parsed:%s", z, m.String())
		}
	}

}

func TestMsgHost(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.host, p.host) {
			t.Errorf("failed:%s parsed:%s", z, m.String())
		}
	}
}
func TestMsgTrailing(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.trailing, p.trailing) {
			t.Errorf("failed:%s parsed:%s", z, m.String())
		}
	}
}
func TestMsgParams(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil {
			for i := 0; i < len(p.params); i++ {
				bytes.Equal(p.params[i], m.params[i])
			}
			t.Errorf("failed:%s parsed:%s", z, m.String())
		}
	}
}
