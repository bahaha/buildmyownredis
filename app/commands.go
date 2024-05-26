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
	Process(*Command, Parser, Storage) ([]byte, error)
}

var (
	handlers = map[string]CommandProcessor{
		"COMMAND": &Great{},
		"PING":    &Ping{},
		"ECHO":    &Echo{},
		"SET":     &Set{},
		"GET":     &Get{},
	}
)

func HandleCommand(cmd *Command, p Parser, s Storage) ([]byte, error) {
	name := strings.ToUpper(string(cmd.Name))
	handler, ok := handlers[name]
	if !ok {
		return nil, errors.New("unknown command: " + name)
	}

	return handler.Process(cmd, p, s)
}

type Great struct{}

func (g *Great) Process(cmd *Command, parser Parser, s Storage) ([]byte, error) {
	return parser.MarshalSimpleString("OK")
}

type Ping struct{}

func (p *Ping) Process(cmd *Command, parser Parser, s Storage) ([]byte, error) {
	return parser.MarshalSimpleString("PONG")
}

type Echo struct{}

func (e *Echo) Process(cmd *Command, p Parser, s Storage) ([]byte, error) {
	msg := cmd.Args[0]
	return p.MarshalBulkString(string(msg))
}

type Set struct{}

func (s *Set) Process(cmd *Command, p Parser, storage Storage) ([]byte, error) {
	key := cmd.Args[0]
	value := cmd.Args[1]

	err := storage.Set(key, value)
	if err != nil {
		return nil, err
	}

	return p.MarshalSimpleString("OK")
}

type Get struct{}

func (g *Get) Process(cmd *Command, p Parser, storage Storage) ([]byte, error) {
	key := cmd.Args[0]

	value, err := storage.Get(key)
	if err != nil {
		return nil, err
	}

	return p.MarshalBulkString(string(value))
}
