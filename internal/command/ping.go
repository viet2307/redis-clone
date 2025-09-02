package command

import (
	"errors"

	"tcp-server.com/m/internal/protocol"
)

func CmdPING(args []string) []byte {
	var res []byte

	if len(args) > 1 {
		return protocol.Encoder(errors.New("ERR wrong number of arguments for 'ping'"), false)
	}
	if len(args) == 0 {
		res = protocol.Encoder("PONG", true)
		return res
	}
	res = protocol.Encoder(args[0], false)
	return res
}
