package datastructure

import (
	"strconv"
	"sync"
)

type Storage struct {
	mu        sync.RWMutex
	dict      Dict
	sortedSet map[string]*ZSet
	cms       map[string]*CMS
	bf        map[string]*Bloom
}

func NewStorage() *Storage {
	return &Storage{
		dict: Dict{
			dictStore:        make(map[string]*Obj),
			expiredDictStore: make(map[string]uint64),
		},
		sortedSet: make(map[string]*ZSet),
	}
}

func (s *Storage) NewCMS(key string, errRate float64, errProb float64) int {
	if _, ok := s.cms[key]; ok {
		return -1
	}
	s.cms[key] = NewCMS(errRate, errProb)
	return 1
}

func (s *Storage) NewBF(key string, errRate float64, entriesNum uint64) int {
	if _, ok := s.bf[key]; ok {
		return -1
	}
	s.bf[key] = NewBloom(errRate, entriesNum)
	return 1
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

/*
Create new zdict for `key`
*/
func (s *Storage) zdictExisted(key string) {
	if _, ok := s.sortedSet[key]; !ok {
		s.sortedSet[key] = NewZset()
	}
}

/*
Currently only support single `element - score` zadd
*/
func (s *Storage) Zadd(key string, args []string) int {
	if len(args) != 2 {
		return -1
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.zdictExisted(key)

	ele := args[0]
	score, _ := strconv.ParseFloat(args[1], 64)
	return s.sortedSet[key].Zadd(ele, score)
}

func (s *Storage) Zscore(key string, ele string) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.zdictExisted(key)

	res := s.sortedSet[key].Zscore(ele)
	if res == nil {
		return -1
	}
	return res.(float64)
}

func (s *Storage) Zrank(key string, ele string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.zdictExisted(key)

	return s.sortedSet[key].Zrank(ele)
}

func (s *Storage) CMSIncrBy(key string, item string, value uint32) uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cms[key].IncrBy(item, value)
}

func (s *Storage) CMSQuery(key string, item string) uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cms[key].Query(item)
}

func (s *Storage) BFAdd(key string, item string) int {
	if _, ok := s.bf[key]; !ok {
		return -1
	}
	s.bf[key].Add(item)
	return 1
}

func (s *Storage) BFQuery(key string, item string) int {
	if _, ok := s.bf[key]; !ok {
		return -1
	}
	res := s.bf[key].Exist(item)
	if res {
		return 1
	}

	return 0
}
