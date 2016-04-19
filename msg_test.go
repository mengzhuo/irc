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
	{
		":WiZ!jto@tolsun.oulu.fi MODE #eu-opers -l",
		&Msg{cmd: s2b("MODE"),
			trailing: nil,
			name:     s2b("WiZ"),
			user:     s2b("jto"),
			host:     s2b("tolsun.oulu.fi"),
			params:   [][]byte{s2b("#eu-opers"), s2b("-l")},
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
			params: [][]byte{s2b("&oulu"),
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
			params:   [][]byte{s2b("#channel")},
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
			params:   [][]byte{s2b("test")},
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
			params:   [][]byte{s2b("test")},
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
			params:   [][]byte{s2b("test")},
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
			params:   [][]byte{s2b("#viveslan")},
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
			params:   [][]byte{s2b("#viveslan")},
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
			params:   [][]byte{},
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
			params:   [][]byte{},
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
			params:   [][]byte{},
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
			params:   [][]byte{s2b("$@")},
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
			params:   [][]byte{s2b("$@"), s2b("param")},
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
			params:   [][]byte{s2b("#foo")},
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
			params:   [][]byte{s2b("#test")},
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
