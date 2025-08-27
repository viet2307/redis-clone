package server

import (
	"bufio"
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

func (h *Handler) HandleConnection() {
	defer func() {
		_ = h.conn.Close()
		log.Printf("Client disconnected: %s\n", h.conn.RemoteAddr().String())
	}()
	fmt.Fprintln(h.conn, "Welcome to the TCP Server! Send 'quit', 'bye' or 'exit' to disconnect.")
	command := map[string]struct{}{
		"quit": {},
		"bye":  {},
		"exit": {},
	}

	r := bufio.NewScanner(h.conn)
	for r.Scan() {
		nextLine := r.Text()
		log.Println("Read line: ", nextLine)
		if _, ok := command[nextLine]; ok {
			fmt.Fprintf(h.conn, "Goodbye!!!")
			return
		}
		fmt.Fprintln(h.conn, "You said:", nextLine)
	}
}
