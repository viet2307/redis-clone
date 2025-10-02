package datastructure

import "math/rand"

type SkiplistNode struct {
	ele      string
	score    float64
	backward *SkiplistNode
	levels   []SkiplistLevel
}

type SkiplistLevel struct {
	forward *SkiplistNode
	span    uint32
}

type Skiplist struct {
	head     *SkiplistNode
	tail     *SkiplistNode
	length   uint32
	level    int
	maxLevel int
}

type ZSet struct {
	zskiplist *Skiplist
	dict      map[string]float64
}

func NewSkiplist(maxLevel int, p float64) *Skiplist {
	head := &SkiplistNode{
		ele:      "",
		backward: nil,
		score:    0,
		levels:   make([]SkiplistLevel, maxLevel),
	}
	return &Skiplist{
		head:     head,
		tail:     nil,
		level:    1,
		length:   0,
		maxLevel: maxLevel,
	}
}

func (s *Skiplist) coinFlip() int {
	x, h := rand.Uint64(), 1
	for x&1 == 1 && h < s.maxLevel {
		h++
		x >>= 1
	}
	return h
}

func (s *Skiplist) newNode(ele string, score float64, h int) *SkiplistNode {
	levels := make([]SkiplistLevel, h)
	return &SkiplistNode{
		ele:    ele,
		score:  score,
		levels: levels,
	}
}

func (s *Skiplist) getBackList(node *SkiplistNode) []*SkiplistNode {
	backList := make([]*SkiplistNode, s.level)
	curr := s.head
	for i := s.level - 1; i >= 0; i-- {
		next := curr.levels[i].forward
		for next != nil && (next.score < node.score || (next.score == node.score && next.ele < node.ele)) {
			curr = next
			next = curr.levels[i].forward
		}
		backList[i] = curr
	}
	return backList
}

func (s *Skiplist) del(node *SkiplistNode, backList []*SkiplistNode) {
	if backList == nil {
		backList = s.getBackList(node)
	}
	for i := 0; i < len(node.levels); i++ {
		back := backList[i]
		next := node.levels[i].forward
		back.levels[i].forward = next
		if i == 0 {
			if next != nil {
				next.backward = back
			} else {
				s.tail = back
			}
		}
	}
	for s.level > 1 && s.head.levels[s.level-1].forward == nil {
		s.level--
	}
	s.length--
	if node == s.tail {
		if node.backward == s.head {
			s.tail = nil
		} else {
			s.tail = node.backward
		}
	}
}

func (s *Skiplist) Zadd(key string, score float64) (interface{}, bool) {
	h := s.coinFlip()
	node := s.newNode(key, score, h)
	backList := s.getBackList(node)
	if h > s.level {
		for i := s.level; i < h; i++ {
			backList = append(backList, s.head)
		}
		s.level = h
	}
	if backList[0].levels[0].forward != nil &&
		backList[0].levels[0].forward.score == score &&
		backList[0].levels[0].forward.ele == key {
		return nil, false
	}

	for i := 0; i < h; i++ {
		node.levels[i].forward = backList[i].levels[i].forward
		backList[i].levels[i].forward = node
	}

	node.backward = backList[0]
	if node.levels[0].forward != nil {
		node.levels[0].forward.backward = node
	} else {
		s.tail = node
	}

	s.length++
	return nil, false
}
