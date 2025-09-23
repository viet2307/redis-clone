package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Obj struct {
	Value interface{}
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

type Server struct {
	listener net.Listener
	port     string
	dict     Dict
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
		dict: Dict{},
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

func (s *Server) Set(key string, value interface{}, expir uint64) {

	if s.dict.dictStore == nil {
		s.dict.dictStore = make(map[string]*Obj)
	}
	if s.dict.expiredDictStore == nil {
		s.dict.expiredDictStore = make(map[string]uint64)
	}
	s.dict.dictStore[key] = &Obj{Value: value}
	s.dict.expiredDictStore[key] = expir
}

func (s *Server) Get(key string) (Obj, bool) {
	if s.dict.dictStore == nil {
		return Obj{}, false
	}
	obj, exist := s.dict.dictStore[key]
	if !exist {
		return Obj{}, false
	}

	if expiredAt, hasExpired := s.dict.expiredDictStore[key]; hasExpired {
		if uint64(time.Now().UnixMilli()) > expiredAt {
			delete(s.dict.dictStore, key)
			delete(s.dict.expiredDictStore, key)
			return Obj{}, false
		}
	}
	return *obj, true
}

func (s *Server) Ttl(key string) (uint64, bool) {
	if s.dict.expiredDictStore == nil {
		return 0, false
	}
	expir, exist := s.dict.expiredDictStore[key]
	if !exist {
		return 0, false
	}
	return expir, true
}

func (s *Server) Expire(key string, expr uint64) (int, bool) {
	if s.dict.expiredDictStore == nil {
		return 0, false
	}
	_, ok := s.dict.expiredDictStore[key]
	if !ok {
		return 0, false
	}
	s.dict.expiredDictStore[key] = expr
	return 1, true
}

func (s *Server) Del(keys []string) (int, bool) {
	if s.dict.dictStore == nil {
		return 0, false
	}
	cnt := 0
	for _, k := range keys {
		if _, ok := s.dict.dictStore[k]; !ok {
			continue
		}
		delete(s.dict.dictStore, k)
		if _, hasExpr := s.dict.expiredDictStore[k]; hasExpr {
			delete(s.dict.expiredDictStore, k)
		}
		cnt++
	}
	return cnt, true
}

func (s *Server) Exist(keys []string) (int, bool) {
	if s.dict.dictStore == nil {
		return 0, false
	}
	cnt := 0
	expiredList := make([]string, 0)
	for _, k := range keys {
		if _, ok := s.dict.dictStore[k]; !ok {
			continue
		}
		if hadExpr, ok := s.dict.expiredDictStore[k]; ok && hadExpr < uint64(time.Now().UnixMilli()) {
			expiredList = append(expiredList, k)
			continue
		}
		cnt++
	}
	s.Del(expiredList)
	return cnt, true
}
