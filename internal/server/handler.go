package server

import (
	"fmt"
	"log"
	"net"
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
	buf := make([]byte, 4096)
	_, err := h.conn.Read(buf)
	if err != nil {
		return
	}

	if _, ok := shutdown[string(buf)]; ok {
		fmt.Fprintf(h.conn, "Goodbye!!!")
		return
	}
	cmd, err := s.executor.CmdParser(buf)
	if err != nil {
		fmt.Fprintf(h.conn, "%s", err)
	}
	h.conn.Write(s.executor.Execute(cmd))
}
