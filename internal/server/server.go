package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"tcp-server.com/m/internal/command"
	"tcp-server.com/m/internal/datastructure"
)

type Server struct {
	listener net.Listener
	port     string
	executor command.Executor
}

func NewServer(port string) *Server {
	return &Server{
		port:     port,
		executor: *command.NewExecutor(datastructure.NewStorage()),
	}
}

func Stop(s *Server) error {
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("error establising new TCP server:\n%v", err)
	}
	s.listener = listen

	log.Printf("Listening on port %s", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// Listener closed via Stop(); exit accept loop
				return nil
			}
			log.Printf("error establishing connection on port %s\n%v", s.port, err)
			continue
		}
		handler := NewHandler(conn)

		go handler.HandleConnection(s)
	}
}
