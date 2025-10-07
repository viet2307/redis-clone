package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tcp-server.com/m/internal/datastructure"
	"tcp-server.com/m/internal/protocol"
)

type Executor struct {
	store *datastructure.Storage
}

type Command struct {
	Name string
	Args []string
}

func NewExecutor(store *datastructure.Storage) *Executor {
	return &Executor{
		store: store,
	}
}

func (e *Executor) CmdParser(data []byte) (*Command, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	in := strings.TrimSpace(string(data))
	cmd, idx := "", 0
	args := make([]string, 0)

	for i, c := range in {
		if c == ' ' && cmd == "" {
			cmd = strings.ToUpper(in[0:i])
			idx = i + 1
			break
		}
	}

	for i := idx; i < len(in); i++ {
		c := in[i]
		if c == ' ' {
			val := strings.ToUpper(in[idx:i])
			args = append(args, val)
			idx = i + 1
		}
	}

	args = append(args, strings.ToUpper(in[idx:]))
	return &Command{
		Name: cmd,
		Args: args,
	}, nil
}

const (
	CmdPing      = "PING"
	CmdSet       = "SET"
	CmdGet       = "GET"
	CmdTtl       = "TTL"
	CmdDel       = "DEL"
	CmdExist     = "EXIST"
	CmdExpire    = "EXPIRE"
	CmdZadd      = "ZADD"
	CmdZScore    = "ZSCORE"
	CmdZrank     = "ZRANK"
	CmdCMSINIT   = "CMS.INITBYPROB"
	CmdCMSIncrBy = "CMS.INCRBY"
	CmdCMSQuery  = "CMS.QUERY"
)

func (e *Executor) Execute(cmd *Command) []byte {
	en := protocol.Encoder{}
	switch cmd.Name {
	case CmdPing:
		return e.cmdPING(cmd.Args)
	case CmdSet:
		return e.cmdSET(cmd.Args)
	case CmdGet:
		return e.cmdGET(cmd.Args)
	case CmdTtl:
		return e.cmdTTL(cmd.Args)
	case CmdExpire:
		return e.cmdExpr(cmd.Args)
	case CmdExist:
		return e.cmdExist(cmd.Args)
	case CmdDel:
		return e.cmdDel(cmd.Args)
	case CmdZadd:
		return e.CmdZadd(cmd.Args)
	case CmdZrank:
		return e.CmdZrank(cmd.Args)
	case CmdZScore:
		return e.CmdZScore(cmd.Args)
	case CmdCMSINIT:
		return e.CmdInitCMS(cmd.Args)
	case CmdCMSIncrBy:
		return e.CmdIncrBy(cmd.Args)
	case CmdCMSQuery:
		return e.CmdCMSQuery(cmd.Args)
	default:
		return en.Encode(errors.New("ERR unsupported CMD detected"), false)
	}
}

func (e *Executor) cmdPING(args []string) []byte {
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

func (e *Executor) cmdSET(args []string) []byte {
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
	e.store.Set(key, val, expr)
	return en.Encode("OK", true)
}

func (e *Executor) cmdGET(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) > 1 || len(args) < 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}
	key := args[0]
	obj, ok := e.store.Get(key)
	if !ok {
		return en.Encode(nil, false)
	}
	return en.Encode(obj.Value, false)
}

func (e *Executor) cmdTTL(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'TTL' command"), false)
	}
	key := args[0]
	ttl, ok := e.store.Ttl(key)
	if !ok {
		return en.Encode(nil, false)
	}
	return en.Encode(ttl, false)
}

func (e *Executor) cmdExpr(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'EXPIRE' command"), false)
	}
	key := args[0]
	expr, _ := strconv.ParseUint(args[1], 10, 64)
	expr += uint64(time.Now().UnixMilli())
	res, _ := e.store.Expire(key, expr)
	return en.Encode(res, false)
}

func (e *Executor) cmdDel(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'DEL' command"), false)
	}
	res, _ := e.store.Del(args)
	return en.Encode(res, false)
}

func (e *Executor) cmdExist(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 1 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'EXIST' command"), false)
	}
	res, _ := e.store.Exist(args)
	return en.Encode(res, false)
}

func (e *Executor) CmdZadd(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) < 3 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'ZADD' command, currently only support single [key, element, score] ZADD"), false)
	}
	res := e.store.Zadd(args[0], args[1:])
	if res == -1 {
		return en.Encode(errors.New("ERR failed to execute command 'ZADD'"), false)
	}
	return en.Encode(res, false)
}

func (e *Executor) CmdZrank(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'ZRANK' command"), false)
	}
	res := e.store.Zrank(args[0], args[1])
	if res == -1 {
		return en.Encode(errors.New("ERR failed to execute command 'ZRANK'"), false)
	}
	return en.Encode(res, false)
}

func (e *Executor) CmdZScore(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'ZSCORE' command"), false)
	}
	res := e.store.Zscore(args[0], args[1])
	if res == -1 {
		return en.Encode(errors.New("ERR failed to execute command 'ZSCORE'"), false)
	}
	return en.Encode(res, false)
}

func (e *Executor) CmdInitCMS(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 3 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'CMS.INITBYPROB' command"), false)
	}
	errRate, _ := strconv.ParseFloat(args[1], 64)
	errProb, _ := strconv.ParseFloat(args[2], 64)
	res := e.store.NewCMS(args[0], errRate, errProb)
	if res == -1 {
		return en.Encode(errors.New("ERR count min sketch with the same key already existed"), false)
	}
	return en.Encode("OK", true)
}

func (e *Executor) CmdIncrBy(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 3 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'CMS.INCRBY' command"), false)
	}
	value, _ := strconv.ParseUint(args[2], 10, 32)
	res := e.store.CMSIncrBy(args[0], args[1], uint32(value))
	return en.Encode(res, false)
}

func (e *Executor) CmdCMSQuery(args []string) []byte {
	en := protocol.Encoder{}
	if len(args) != 2 {
		return en.Encode(errors.New("ERR wrong number of arguments for 'CMS.QUERY' command"), false)
	}
	res := e.store.CMSQuery(args[0], args[1])
	return en.Encode(res, false)
}
