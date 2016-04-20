package irc

import (
	"bytes"
	"strings"
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
			params:   [16][]byte{s2b("cm22.eng.umd.edu")},
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
			params:   [16][]byte{s2b("#playzone")},
		},
		true,
		false,
	},
	{
		":WiZ!jto@tolsun.oulu.fi MODE #eu-opers -l",
		&Msg{cmd: s2b("MODE"),
			trailing: nil,
			name:     s2b("WiZ"),
			user:     s2b("jto"),
			host:     s2b("tolsun.oulu.fi"),
			params:   [16][]byte{s2b("#eu-opers"), s2b("-l")},
		},
		true,
		false,
	},
	{
		"MODE &oulu +b *!*@*.edu +e *!*@*.bu.edu",
		&Msg{cmd: s2b("MODE"),
			trailing: nil,
			name:     nil,
			user:     nil,
			host:     nil,
			params: [16][]byte{s2b("&oulu"),
				s2b("+b"),
				s2b("*!*@*.edu"),
				s2b("+e"),
				s2b("*!*@*.bu.edu"),
			},
		},
		false,
		true,
	},
	{
		"PRIVMSG #channel :Message with :colons!",
		&Msg{cmd: s2b("PRIVMSG"),
			trailing: s2b("Message with :colons!"),
			name:     nil,
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("#channel")},
		},
		false,
		true, // TODO PRIVMSG should not be server
	},
	{
		":irc.vives.lan 251 test :There are 2 users and 0 services on 1 servers",
		&Msg{cmd: s2b("251"),
			trailing: s2b("There are 2 users and 0 services on 1 servers"),
			name:     s2b("irc.vives.lan"),
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("test")},
		},
		false,
		true,
	},
	{
		":irc.vives.lan 376 test :End of MOTD command",
		&Msg{cmd: s2b("376"),
			trailing: s2b("End of MOTD command"),
			name:     s2b("irc.vives.lan"),
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("test")},
		},
		false,
		true,
	},
	{
		":irc.vives.lan 250 test :Highest connection count: 1 (1 connections received)",
		&Msg{cmd: s2b("250"),
			trailing: s2b("Highest connection count: 1 (1 connections received)"),
			name:     s2b("irc.vives.lan"),
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("test")},
		},
		false,
		true,
	},
	{
		":sorcix!~sorcix@sorcix.users.quakenet.org PRIVMSG #viveslan :\001ACTION is testing CTCP messages!\001",
		&Msg{cmd: s2b("PRIVMSG"),
			trailing: s2b("\001ACTION is testing CTCP messages!\001"),
			name:     s2b("sorcix"),
			user:     s2b("~sorcix"),
			host:     s2b("sorcix.users.quakenet.org"),
			params:   [16][]byte{s2b("#viveslan")},
		},
		true,
		false,
	},
	{
		":sorcix!~sorcix@sorcix.users.quakenet.org NOTICE midnightfox :\001PONG 1234567890\001",
		&Msg{cmd: s2b("NOTICE"),
			trailing: s2b("\001PONG 1234567890\001"),
			name:     s2b("sorcix"),
			user:     s2b("~sorcix"),
			host:     s2b("sorcix.users.quakenet.org"),
			params:   [16][]byte{s2b("midnightfox")},
		},
		true,
		false,
	},
	{
		":a!b@c QUIT",
		&Msg{cmd: s2b("QUIT"),
			trailing: nil,
			name:     s2b("a"),
			user:     s2b("b"),
			host:     s2b("c"),
			params:   [16][]byte{},
		},
		true,
		false,
	},
	{
		":a!b PRIVMSG :message",
		&Msg{cmd: s2b("PRIVMSG"),
			trailing: s2b("message"),
			name:     s2b("a"),
			user:     s2b("b"),
			host:     nil,
			params:   [16][]byte{},
		},
		false,
		false,
	},
	{
		":a@c NOTICE ::::Hey!",
		&Msg{cmd: s2b("NOTICE"),
			trailing: s2b(":::Hey!"),
			name:     s2b("a"),
			user:     nil,
			host:     s2b("c"),
			params:   [16][]byte{},
		},
		false,
		false,
	},
	{
		":nick PRIVMSG $@ :This message contains a\ttab!",
		&Msg{cmd: s2b("PRIVMSG"),
			trailing: s2b("This message contains a\ttab!"),
			name:     s2b("nick"),
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("$@")},
		},
		false,
		true, // TODO
	},
	{
		"TEST $@  param :Trailing",
		&Msg{cmd: s2b("TEST"),
			trailing: s2b("Trailing"),
			name:     nil,
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("$@"), s2b("param")},
		},
		false,
		true,
	},
	{
		"TOPIC #foo",
		&Msg{cmd: s2b("TOPIC"),
			trailing: nil,
			name:     nil,
			user:     nil,
			host:     nil,
			params:   [16][]byte{s2b("#foo")},
		},
		false,
		true,
	},
	{
		":name!user@example.org PRIVMSG #test :Message with spaces at the end!  ",
		&Msg{cmd: s2b("PRIVMSG"),
			trailing: s2b("Message with spaces at the end!  "),
			name:     s2b("name"),
			user:     s2b("user"),
			host:     s2b("example.org"),
			params:   [16][]byte{s2b("#test")},
		},
		true,
		false,
	},
}

func TestInvaildMsg(t *testing.T) {
	invalid := []string{
		": PRIVMSG test :Invalid message with empty prefix.",
		":  PRIVMSG test :Invalid message with space prefix",
	}
	for _, s := range invalid {
		m, err := NewMsg(s2b(s))
		if err == nil || m.Cmd() != nil {
			t.Error(s, "is valid")
		}

	}
}

func TestMsgCmd(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.Cmd(), p.cmd) {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}

func TestMsgPrefixName(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.Name(), p.name) {
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
		if err != nil || !bytes.Equal(m.User(), p.user) {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}

}

func TestMsgHost(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.Host(), p.host) {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}
func TestMsgTrailing(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		p := z.parsed
		if err != nil || !bytes.Equal(m.Trailing(), p.trailing) {
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
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
		var j int
		var zz []byte
		for j, zz = range z.parsed.params {
			if zz == nil {
				break
			}
		}
		for i := 0; i < j; i++ {
			if !bytes.Equal(m.Params()[i], p.params[i]) {
				t.Log(z, "i=", i, "j=", j, m.Params()[i], p.params[i])
				t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
			}
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

func TestMsgString(t *testing.T) {
	for _, z := range messageTests {
		m, err := NewMsg(s2b(z.rawMsg))
		m.ParseAll()
		if err != nil || strings.HasPrefix(m.String(), "CMD: ") {
			t.Errorf("failed:%s\nparsed:%s", z.rawMsg, m.String())
		}
	}
}

func TestSetCmd(t *testing.T) {
	m := new(Msg)
	m.SetCmd(s2b("XXX"))
	if string(m.Cmd()) != "XXX" {
		t.Error(m.Cmd(), "!=", "XXX")
	}
}

func TestSetHost(t *testing.T) {
	m := new(Msg)
	m.SetHost(s2b("XXX"))
	if string(m.host) != "XXX" {
		t.Error(m.host, "!=", "XXX")
	}
}

func TestSetUser(t *testing.T) {
	m := new(Msg)
	m.SetUser(s2b("XXX"))
	if string(m.user) != "XXX" {
		t.Error(m.user, "!=", "XXX")
	}
}

func TestSetName(t *testing.T) {
	m := new(Msg)
	m.SetName(s2b("XXX"))
	if string(m.name) != "XXX" {
		t.Error(m.name, "!=", "XXX")
	}
}
func TestSetTrailing(t *testing.T) {
	m := new(Msg)
	m.SetTrailing(s2b("XXX"))
	if string(m.Trailing()) != "XXX" {
		t.Error(m.trailing, "!=", "XXX")
	}
}

func BenchmarkParseMessage_short(b *testing.B) {
	src := s2b("COMMAND arg1 :Message\r\n")
	m := new(Msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Data = src
		m.ParseAll()
		m.Reset()
	}
}

func BenchmarkParseMessage_medium(b *testing.B) {
	src := s2b(":Namename COMMAND arg6 arg7 :Message message message\r\n")
	m := new(Msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Data = src
		m.ParseAll()
		m.Reset()
	}
}
func BenchmarkParseMessage_long(b *testing.B) {
	src := s2b(":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n")
	m := new(Msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Data = src
		m.ParseAll()
		m.Reset()
	}
}

var params [][]byte

func BenchmarkParamsAlloc(b *testing.B) {
	src := s2b(":Namename!username@hostname COMMAND arg1 arg2 arg3 arg4 arg5 arg6 arg7 :Message message message message message\r\n")
	m := new(Msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Data = src
		m.PeekCmd()
		params = m.Params()
		m.Reset()
	}
}

func TestAppendParams(t *testing.T) {
	m := new(Msg)
	m.AppendParams(s2b("#ERR"))
	if m.paramsCount != 1 || string(m.params[0]) != "#ERR" {
		t.Error(m)
	}
}

func TestAppendParamsOverflow(t *testing.T) {
	n := bytes.Split(bytes.Repeat([]byte(" param"), 16)[1:], []byte{space})
	m := new(Msg)
	m.prefixParsed = true
	m.SetCmd([]byte("hehe"))
	if err := m.SetParams(n...); err != nil || m.paramsCount != 16 {
		t.Error(err, n, m)
	}
	if err := m.AppendParams([]byte("$")); err == nil {
		t.Error(err, m, n)
	}
}

func TestSetParams(t *testing.T) {
	m := new(Msg)
	m.SetParams(s2b("1"), s2b("2"))
	if m.paramsCount != 2 || string(m.params[0]) != "1" || string(m.params[1]) != "2" {
		t.Error(m)
	}

}

func TestSetParamsOver(t *testing.T) {
	n := bytes.Split(bytes.Repeat([]byte("param "), 17), []byte{space})
	m := new(Msg)
	m.prefixParsed = true
	m.SetCmd([]byte("hehe"))
	err := m.SetParams(n...)
	if err == nil || m.params[0] != nil {
		t.Error(m, n)
	}
}
