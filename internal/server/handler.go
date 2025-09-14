package server

import (
	"fmt"
	"log"
	"net"

	"tcp-server.com/m/internal/command"
	"tcp-server.com/m/internal/protocol"
)

type Handler struct {
	conn net.Conn
}

func NewHandler(conn net.Conn) *Handler {
	return &Handler{
		conn: conn,
	}
}

func (h *Handler) HandleConnection() {
	defer func() {
		_ = h.conn.Close()
		log.Printf("Client disconnected: %s\n", h.conn.RemoteAddr().String())
	}()
	fmt.Fprintln(h.conn, "Welcome to the TCP Server! Send 'quit', 'bye' or 'exit' to disconnect.")
	shutdown := map[string]struct{}{
		"quit": {},
		"bye":  {},
		"exit": {},
	}
	buf := make([]byte, 4096)
	_, err := h.conn.Read(buf)
	if err != nil {
		return
	}

	if _, ok := shutdown[string(buf)]; ok {
		fmt.Fprintf(h.conn, "Goodbye!!!")
		return
	}

	cmd, err := protocol.CmdParser(buf)
	if err != nil {
		fmt.Fprintf(h.conn, "%s", err)
	}

	switch cmd.Cmd {
	case "PING":
		h.conn.Write(command.CmdPING(cmd.Args))
	}
	// cmd, _, err := protocol.ArrParser(buf[:n], 0)
	// if err != nil {
	// 	log.Println("Parse error: ", err)
	// 	return
	// }
	// if len(cmd) == 0 {
	// 	log.Printf("ERR command empty")
	// 	return
	// }
	// cmdName := strings.ToUpper(cmd[0].(string))

	// args := make([]string, 0, len(cmd)-1)
	// for _, arg := range cmd[1:] {
	// 	args = append(args, arg.(string))
	// }

	// switch cmdName {
	// case "PING":
	// 	h.conn.Write(command.CmdPING(args))
	// default:
	// 	h.conn.Write(protocol.Encoder(errors.New("ERR unknown command '"+cmdName+"'"), false))
	// }
}
