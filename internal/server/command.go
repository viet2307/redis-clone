package server

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tcp-server.com/m/internal/protocol"
)

const (
	CmdPing   = "PING"
	CmdSet    = "SET"
	CmdGet    = "GET"
	CmdTtl    = "TTL"
	CmdDel    = "DEL"
	CmdExist  = "EXIST"
	CmdExpire = "EXPIRE"
)

type Command struct {
	Cmd  string
	Args []string
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

func (c *Command) Execute(s *Server) []byte {
	en := protocol.Encoder{}
	switch c.Cmd {
	case CmdPing:
		return cmdPING(c.Args)
	case CmdSet:
		return cmdSET(s, c.Args)
	case CmdGet:
		return cmdGET(s, c.Args)
	case CmdTtl:
		return cmdTTL(s, c.Args)
	case CmdExpire:
		return cmdExpr(s, c.Args)
	case CmdExist:
		return cmdExist(s, c.Args)
	case CmdDel:
		return cmdDel(s, c.Args)
	default:
		return en.Encode(errors.New("ERR unsupported CMD detected"), false)
	}
}

func cmdPING(args []string) []byte {
	var res []byte
	en := protocol.Encoder{}
	if len(args) > 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}
	if len(args) == 0 {
		res = en.Encode("PONG", true)
		return res
	}
	res = en.Encode(args[0], false)
	return res
}

func cmdSET(s *Server, args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'set' command"), false)
	}
	key, val := args[0], args[1]
	var ttl uint64 = 0
	if len(args) == 4 {
		parsed, _ := strconv.ParseInt(args[3], 10, 64)
		ttl = uint64(parsed)
	}
	expr := ttl + uint64(time.Now().UnixMilli())
	s.Set(key, val, expr)
	return en.Encode("OK", true)
}

func cmdGET(s *Server, args []string) []byte {
	en := protocol.Encoder{}
	if len(args) > 1 || len(args) < 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}
	key := args[0]
	obj, ok := s.Get(key)
	if !ok {
		return en.Encode(nil, false)
	}
	return en.Encode(obj.Value, false)
}

func cmdTTL(s *Server, args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'TTL' command"), false)
	}
	key := args[0]
	ttl, ok := s.Ttl(key)
	if !ok {
		return en.Encode(nil, false)
	}
	return en.Encode(ttl, false)
}

func cmdExpr(s *Server, args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'EXPIRE' command"), false)
	}
	key := args[0]
	expr, _ := strconv.ParseUint(args[1], 10, 64)
	expr += uint64(time.Now().UnixMilli())
	res, _ := s.Expire(key, expr)
	return en.Encode(res, false)
}

func cmdDel(s *Server, args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'DEL' command"), false)
	}
	keys := args[1:]
	res, _ := s.Del(keys)
	return en.Encode(res, false)
}

func cmdExist(s *Server, args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'EXIST' command"), false)
	}
	keys := args[1:]
	res, _ := s.Exist(keys)
	return en.Encode(res, false)
}
