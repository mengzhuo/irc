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
	{
		":Trillian SQUIT cm22.eng.umd.edu :Server out of control",
		&Msg{cmd: s2b("SQUIT"),
			trailing: s2b("Server out of control"),
			name:     s2b("Trillian"),
			user:     nil,
			host:     nil,
			params:   [][]byte{s2b("cm22.eng.umd.edu")},
		},
		false,
		true,
	},
	{
		":WiZ!jto@tolsun.oulu.fi PART #playzone :I lost",
		&Msg{cmd: s2b("PART"),
			trailing: s2b("I lost"),
			name:     s2b("WiZ"),
			user:     s2b("jto"),
			host:     s2b("tolsun.oulu.fi"),
			params:   [][]byte{s2b("#playzone")},
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
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}

func TestMsgPrefixName(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.name, p.name) {
			t.Errorf("failed:%s\nparsed:%s\nm=%s p=%s", z.rawMsg, m.String(),
				m.name, p.name)
		}
	}
}

func TestMsgUser(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.user, p.user) {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}

}

func TestMsgHost(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.host, p.host) {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}
func TestMsgTrailing(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.trailing, p.trailing) {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
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
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}

func TestMsgServer(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		if err != nil || m.IsServer() != z.server {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}

func TestMsgHostMask(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		if err != nil || m.IsHostMask() != z.hostmask {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}
