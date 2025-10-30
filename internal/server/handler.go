package server

import (
	"fmt"
	"log"
	"net"
	"strings"

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

func (h *Handler) HandleConnection(s *Server) {
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
	for {
		buf := make([]byte, 4096)
		n, err := h.conn.Read(buf)
		if err != nil {
			return
		}

		data := strings.TrimSpace(string(buf[:n]))
		if _, ok := shutdown[strings.ToLower(string(data))]; ok {
			fmt.Fprintf(h.conn, "Goodbye!!!")
			return
		}
		cmd, err := s.executor.CmdParser([]byte(data))
		if err != nil {
			fmt.Fprintf(h.conn, "%s", err)
		}
		parser := protocol.REPSParser{}
		res, err := parser.Parse(s.executor.Execute(cmd))
		if err != nil {
			fmt.Fprintf(h.conn, "%s", err)
		}
		res = append(res, '\r', '\n')
		h.conn.Write(res)
	}
}
