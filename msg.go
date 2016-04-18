package irc

import (
	"bytes"
	"errors"
	"fmt"
)

type Msg struct {
	Data         []byte
	Prefix       []byte
	Cmd          []byte
	Params       [][]byte
	Trailing     []byte
	parsedCmd    bool
	parsedParams bool
	index        int
}

func NewMsg(b []byte) (m *Msg) {
	m = new(Msg)
	m.Data = b
	return
}

func (m *Msg) parseAll() (err error) {

	err = m.parseCmd()
	if err != nil {
		return
	}
	err = m.parseParams()
	return
}

func (m *Msg) parseCmd() (err error) {
	var n int
	if m.parsedCmd {
		return
	}

	b := m.Data
	if m.hasPrefix() {
		n = bytes.IndexByte(b, ' ')
		if n <= 0 {
			err = errors.New("can't get cmd")
			return
		} else {
			m.Prefix = b[:n]
			m.index = n
			b = b[n+1:]
		}
	}

	n = bytes.IndexByte(b, ' ')
	if n < 0 {
		// no params
		m.Cmd = b
		m.index = m.index + len(b) + 1
	} else {
		m.Cmd = b[:n]
		m.index += n + 1
	}

	m.parsedCmd = true
	return err
}

func (m *Msg) parseParams() (err error) {
	var n int
	if m.parsedParams {
		return
	}
	b := m.Data[m.index:]
	// find trailing
	n = bytes.LastIndexByte(b, ':')
	if n < 0 {
		// No trailing
		n = len(b) - 1
	} else {
		m.Trailing = b[n:]
		b = b[:n]
	}

	for {
		n = bytes.IndexByte(b, ' ')
		if n < 0 {
			break
		}
		m.Params = append(m.Params, b[:n])
		b = b[n+1:]
	}

	m.parsedParams = true
	return
}

func (m *Msg) hasPrefix() bool {
	return m.Data[0] == ':'
}

func (m *Msg) String() string {
	return fmt.Sprintf("CMD:%s, Params:%s, Prefix=%s Trailing=%s",
		m.Cmd, m.Params, m.Prefix, m.Trailing)
}

func (m *Msg) Reset() {
	m.Data = nil
	m.Prefix = nil
	m.Cmd = nil
	m.Params = nil
	m.Trailing = nil
	m.parsedParams = false
	m.parsedCmd = false
}
