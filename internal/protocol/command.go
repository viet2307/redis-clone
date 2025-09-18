package protocol

import (
	"errors"
	"fmt"
	"strings"
)

const (
	CmdPing = "PING"
	CmdSet  = "SET"
	CmdGET  = "GET"
)

type Command struct {
	Cmd  string
	Args []string
}

type Obj struct {
	Value interface{}
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func (c *Command) CmdParser(data []byte) (Command, error) {
	if len(data) == 0 {
		return Command{}, fmt.Errorf("empty input")
	}
	cmd := ""
	idx := 0
	args := make([]string, 0)
	for i, c := range data {
		if c == ' ' && cmd == "" {
			cmd = strings.ToLower(string(data)[0:i])
			idx = i + 1
		} else {
			val := strings.ToLower(string(data)[idx:i])
			args = append(args, val)
			idx = i + 1
		}
	}
	return Command{
		Cmd:  cmd,
		Args: args,
	}, nil
}

func (c *Command) Execute() []byte {
	en := Encoder{}
	switch c.Cmd {
	case CmdPing:
		return cmdPING(c.Args)
	case CmdSet:
		return cmdSET(c.Args)
	default:
		return en.Encode(errors.New("ERR invalid CMD"))
	}
}

func cmdPING(args []string) []byte {
	var res []byte
	en := Encoder{}
	if len(args) > 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'ping'"), false)
	}
	if len(args) == 0 {
		res = en.Encode("PONG", true)
		return res
	}
	res = en.Encode(args[0], false)
	return res
}

func cmdSET(args []string) []byte {
	return nil
}
