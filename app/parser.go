package main

import (
	"bufio"
	"bytes"
	"errors"
	"log/slog"
	"strconv"
)

type Parser interface {
	Parse([]byte) (*Command, error)

	MarshalSimpleString(string) ([]byte, error)
	MarshalBulkString(string) ([]byte, error)
}

var (
	DELIM = []byte("\r\n")
)

type RESP struct {
}

func NewResp() *RESP {
	return &RESP{}
}

func (p *RESP) Parse(raw []byte) (*Command, error) {
	slog.Info("Received command bytes from peer", "msg", raw)
	headByte := raw[0]

	if headByte == '*' {
		return p.readArray(raw)
	}

	return &Command{}, nil
}

func (p *RESP) readArray(raw []byte) (*Command, error) {
	reader := bufio.NewReader(bytes.NewReader(raw))
	if hb, err := reader.ReadByte(); err != nil || hb != '*' {
		return nil, errors.New("invalid RESP format: expected '*' as first byte of an array")
	}

	n, err := p.readInt(reader)
	if err != nil {
		return nil, err
	}

	cmd := &Command{
		Args: make([][]byte, n-1),
	}

	name, err := p.readBulkString(reader)
	if err != nil {
		return nil, err
	}
	cmd.Name = name
	for i := 1; i < n; i++ {
		arg, err := p.readBulkString(reader)
		if err != nil {
			return nil, err
		}

		cmd.Args[i-1] = arg
	}

	return cmd, nil
}

func (p *RESP) readInt(r *bufio.Reader) (int, error) {
	line, err := p.readLine(r)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(line))
}

func (p *RESP) readBulkString(r *bufio.Reader) ([]byte, error) {
	if hb, err := r.ReadByte(); err != nil || hb != '$' {
		return nil, errors.New("invalid RESP format: expected '*' as first byte of an array")
	}

	n, err := p.readInt(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, n+len(DELIM))
	r.Read(buf)

	return bytes.TrimSuffix(buf, DELIM), nil
}

func (p *RESP) readLine(r *bufio.Reader) ([]byte, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(line, DELIM), nil
}

func (p *RESP) MarshalSimpleString(s string) ([]byte, error) {
	return []byte("+" + s + "\r\n"), nil
}

func (p *RESP) MarshalBulkString(s string) ([]byte, error) {
	return []byte("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"), nil
}
