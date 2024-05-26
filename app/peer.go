package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	conn       net.Conn
	remoteAddr net.Addr
	parser     Parser
	storage    Storage
}

func NewPeer(conn net.Conn, parser Parser, storage Storage) *Peer {
	return &Peer{
		conn:       conn,
		remoteAddr: conn.RemoteAddr(),
		parser:     parser,
		storage:    storage,
	}
}

func (p *Peer) Close() error {
	return p.conn.Close()
}

func (p *Peer) WaitForCommand() {
	for {
		buf := make([]byte, 1024)
		nBytes, err := p.conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Connection closed by peer: ", "remoteAddr", p.remoteAddr)
				return
			}
			slog.Error("Error reading from peer: ", "error", err)
		}

		raw := make([]byte, nBytes)
		copy(raw, buf[:nBytes])
		cmd, err := p.parser.Parse(raw)
		if err != nil {
			slog.Error("Error parsing command:", "error", err)
			return
		}
		slog.Info("Received command from peer", "cmd", cmd)

		resp, err := HandleCommand(cmd, p.parser, p.storage)
		if err != nil {
			slog.Error("Error handling command:", "error", err)
			return
		}

		slog.Info("Sending response to peer", "response", string(resp))
		_, err = p.conn.Write(resp)
		if err != nil {
			slog.Error("Error writing to peer:", "error", err)
		}
	}
}
