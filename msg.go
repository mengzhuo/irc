package irc

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	prefixSymbol byte = 0x3a // Prefix Symbol
	userSymbol   byte = 0x21 // Username
	hostSymbol   byte = 0x40 // Hostname
	space        byte = 0x20 // Sepector
)

type Msg struct {
	Data     []byte
	prefix   []byte
	name     []byte
	user     []byte
	host     []byte
	cmd      []byte
	params   [][]byte
	trailing []byte

	prefixParsed bool
	paramsParsed bool
	index        int
}

// Prefix
func (m *Msg) Name() []byte {
	m.parsePrefix()
	return m.name
}

func (m *Msg) User() []byte {
	m.parsePrefix()
	return m.user
}

func (m *Msg) Host() []byte {
	m.parsePrefix()
	return m.host
}

func (m *Msg) IsHostMask() bool {
	m.parsePrefix()
	return m.host != nil && m.user != nil
}

func (m *Msg) IsServer() bool {
	m.parsePrefix()
	return m.user == nil && m.host == nil
}

// Params

func (m *Msg) Params() [][]byte {
	m.parseParams()
	return m.params
}

func (m *Msg) Trailing() []byte {
	m.parseParams()
	return m.trailing
}

func (m *Msg) Cmd() []byte {
	return m.cmd
}

func NewMsg(b []byte) (m *Msg, err error) {
	m = new(Msg)
	m.Data = b
	err = m.PeekCmd()
	return
}

func (m *Msg) PeekCmd() (err error) {

	b := m.Data
	var n int

	if m.hasPrefix() {
		n = bytes.IndexByte(b, space)
		if n <= 0 {
			err = errors.New("can't get cmd")
			return
		} else {
			m.prefix = b[1:n]
			m.index = n
			b = b[n+1:]
		}
	}

	n = bytes.IndexByte(b, space)
	if n < 0 {
		// no params
		m.cmd = b
		m.index = m.index + len(b) + 1
	} else {
		m.cmd = b[:n]
		m.index += n + 1
	}

	return err
}

func (m *Msg) parsePrefix() (err error) {

	if m.prefixParsed {
		return
	}
	user := bytes.IndexByte(m.prefix, userSymbol)
	host := bytes.IndexByte(m.prefix, hostSymbol)

	switch {
	case user > 0 && host > user:
		m.name = m.prefix[:user]
		m.user = m.prefix[user+1 : host]
		m.host = m.prefix[host+1:]

	case user > 0:
		m.name = m.prefix[:user]
		m.user = m.prefix[user+1:]

	case host > 0:
		m.name = m.prefix[:host]
		m.host = m.prefix[host+1:]
	}

	m.prefixParsed = true
	return
}

func (m *Msg) parseParams() (err error) {

	if m.paramsParsed {
		return
	}

	var n int
	b := m.Data[m.index:]
	// find trailing
	n = bytes.LastIndexByte(b, prefixSymbol)
	if n < 0 {
		// No trailing
		n = len(b) - 1
	} else {
		m.trailing = b[n+1:]
		if b[0] == space {
			b = b[1:n]
		} else {
			b = b[:n]
		}
	}

	for {
		n = bytes.IndexByte(b, space)
		if n < 0 {
			break
		}
		m.params = append(m.params, b[:n])
		b = b[n+1:]
	}

	m.paramsParsed = true
	return
}

func (m *Msg) hasPrefix() bool {
	return m.Data[0] == prefixSymbol
}

func (m *Msg) String() string {
	return fmt.Sprintf("CMD:%s, Params:%s, Prefix=%s Trailing=%s",
		m.cmd, m.params, m.prefix, m.trailing)
}

func (m *Msg) Reset() {
	m.Data = nil
	m.prefix = nil
	m.cmd = nil
	m.params = nil
	m.trailing = nil
	m.paramsParsed = false
	m.prefixParsed = false
}
