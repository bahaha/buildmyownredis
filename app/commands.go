package main

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	Name []byte
	Args [][]byte
}

func (c *Command) String() string {
	args := make([]string, len(c.Args))
	for i, arg := range c.Args {
		args[i] = string(arg)
	}

	return fmt.Sprintf(
		"Command{Name=%v, Args=%v}",
		strings.ToUpper(string(c.Name)),
		args,
	)
}

type CommandProcessor interface {
	Process(*Command, Parser) ([]byte, error)
}

var (
	handlers = map[string]CommandProcessor{
		"COMMAND": &Great{},
		"PING":    &Ping{},
		"ECHO":    &Echo{},
	}
)

func HandleCommand(cmd *Command, p Parser) ([]byte, error) {
	handler, ok := handlers[strings.ToUpper(string(cmd.Name))]
	if !ok {
		return nil, errors.New("unknown command")
	}

	return handler.Process(cmd, p)
}

type Great struct{}

func (g *Great) Process(cmd *Command, parser Parser) ([]byte, error) {
	return []byte("+OK\r\n"), nil
}

type Ping struct{}

func (p *Ping) Process(cmd *Command, parser Parser) ([]byte, error) {
	return []byte("+PONG\r\n"), nil
}

type Echo struct{}

func (e *Echo) Process(cmd *Command, p Parser) ([]byte, error) {
	msg := cmd.Args[0]
	return p.MarshalBulkString(string(msg))
}
