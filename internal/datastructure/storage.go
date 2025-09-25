package datastructure

import (
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	dict Dict
}

func NewStorage() *Storage {
	return &Storage{
		dict: Dict{
			dictStore:        make(map[string]*Obj),
			expiredDictStore: make(map[string]uint64),
		},
	}
}

func (s *Storage) Set(key string, value interface{}, expir uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dict.Set(key, value, expir)
}

func (s *Storage) Get(key string) (Obj, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dict.Get(key)
}

func (s *Storage) Ttl(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dict.Ttl(key)
}

func (s *Storage) Expire(key string, expr uint64) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dict.Expire(key, expr)
}

func (s *Storage) Del(keys []string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dict.Del(keys)
}

func (s *Storage) Exist(keys []string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dict.Exist(keys)
}
