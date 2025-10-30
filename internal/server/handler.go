package server

import (
	"fmt"
	"log"
	"net"

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
	for {
		buf := make([]byte, 4096)
		n, err := h.conn.Read(buf)
		if err != nil {
			return
		}

		parser := protocol.REPSParser{}
		cmdParts, err := parser.Parse(buf[:n])
		if err != nil {
			fmt.Fprintf(h.conn, "-ERR %s\r\n", err)
			continue
		}

		cmd, err := s.executor.CmdParser(cmdParts)
		if err != nil {
			fmt.Fprintf(h.conn, "-ERR %s\r\n", err)
			continue
		}

		res := s.executor.Execute(cmd)
		h.conn.Write(res)
	}
}
