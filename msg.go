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
	params   [16][]byte
	trailing []byte

	prefixParsed bool
	paramsParsed bool
	index        int
	paramsCount  int
}

// Prefix

func (m *Msg) makePrefix() {
	if m.name != nil {
		m.prefix = m.name
	}
	if m.user != nil {
		m.prefix = append(m.prefix, userSymbol)
		m.prefix = append(m.prefix, m.user...)
	}
	if m.host != nil {
		m.prefix = append(m.prefix, hostSymbol)
		m.prefix = append(m.prefix, m.host...)
	}
}

func (m *Msg) Name() []byte {
	m.parsePrefix()
	return m.name
}

func (m *Msg) SetName(p []byte) {
	m.name = p
	m.makePrefix()
}

func (m *Msg) User() []byte {
	m.parsePrefix()
	return m.user
}

func (m *Msg) SetUser(p []byte) {
	m.user = p
	m.makePrefix()
}

func (m *Msg) Host() []byte {
	m.parsePrefix()
	return m.host
}

func (m *Msg) SetHost(p []byte) {
	m.host = p
	m.makePrefix()
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
	return m.params[:m.paramsCount]
}

func (m *Msg) SetParams(params ...[]byte) (err error) {

	if len(params) > 16 {
		err = errors.New("Too many params")
		return
	}

	for i, p := range params {
		m.params[i] = p
	}
	m.paramsCount = len(params)
	m.paramsParsed = true

	return
}

func (m *Msg) AppendParams(p []byte) (err error) {

	if m.paramsCount >= 16 {
		err = errors.New("Too many params")
		return
	}

	m.params[m.paramsCount] = p
	m.paramsCount += 1
	m.paramsParsed = true
	return
}

func (m *Msg) Trailing() []byte {
	m.parseParams()
	return m.trailing
}

func (m *Msg) SetTrailing(p []byte) {
	m.trailing = p
	m.paramsParsed = true
}

func (m *Msg) Cmd() []byte {
	return m.cmd
}

func (m *Msg) SetCmd(p []byte) {
	m.cmd = p
}

func NewMsg(b []byte) (m *Msg, err error) {
	m = new(Msg)
	m.Data = b
	err = m.PeekCmd()
	return
}

func (m *Msg) PeekCmd() (err error) {

	if m.index != 0 {
		return
	}

	var n int
	b := m.Data

	if m.hasPrefix() {
		n = bytes.IndexByte(b, space)
		if n == 1 {
			err = errors.New("prefix is empty")
			return

		}

		m.prefix = b[1:n]
		m.index = n
		b = b[n+1:]

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
	default:
		m.name = m.prefix
	}

	m.prefixParsed = true
	return
}

func (m *Msg) parseParams() (err error) {

	if m.paramsParsed {
		return
	}

	var n int
	m.paramsCount = 0
	b := m.Data[m.index:]
	// find trailing
	n = bytes.IndexByte(b, prefixSymbol)
	if n < 0 {
		// No trailing
		n = len(b) - 1
	} else {
		m.trailing = b[n+1:]
		b = b[:n]
	}
	for {
		if len(b) == 0 {
			break
		}
		n = bytes.IndexByte(b, space)
		if n == 0 {
			b = b[1:]
			continue
		}
		if n < 0 {
			m.params[m.paramsCount] = b[:]
			m.paramsCount += 1
			break
		} else {
			m.params[m.paramsCount] = b[:n]
			b = b[n:]
		}
		m.paramsCount += 1
	}

	m.paramsParsed = true
	return
}

func (m *Msg) ParseAll() (err error) {
	err = m.PeekCmd()
	err = m.parsePrefix()
	err = m.parseParams()
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
	m.name = nil
	m.user = nil
	m.host = nil
	m.trailing = nil
	m.paramsParsed = false
	m.prefixParsed = false

	for i := 0; i < m.paramsCount; i++ {
		m.params[i] = nil
	}
	m.paramsCount = 0
	m.index = 0
}
