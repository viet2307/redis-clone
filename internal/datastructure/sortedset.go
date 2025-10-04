package datastructure

import (
	"math/rand"
)

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

func NewSkiplist(maxLevel int) *Skiplist {
	skiplistLevel := make([]SkiplistLevel, maxLevel)
	for i := range len(skiplistLevel) {
		skiplistLevel[i].forward = nil
		skiplistLevel[i].span = 0
	}

	head := &SkiplistNode{
		ele:      "",
		backward: nil,
		score:    0,
		levels:   skiplistLevel,
	}
	return &Skiplist{
		head:     head,
		tail:     nil,
		level:    1,
		length:   0,
		maxLevel: maxLevel,
	}
}

/*
Flip coin to promote new node with probability of 0.5
*/
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

func (s *Skiplist) getBackList(node *SkiplistNode) ([]*SkiplistNode, []uint32) {
	backList := make([]*SkiplistNode, s.level)
	rank := make([]uint32, s.level)
	curr := s.head
	span := uint32(1)
	for i := 0; i < s.level; i++ {
		next := curr.levels[i].forward
		for next != nil && (next.score < node.score || (next.score == node.score && next.ele < node.ele)) {
			curr = next
			next = curr.levels[i].forward
			if i == 0 {
				span++
			}
		}
		backList[i] = curr
		if i == 0 {
			rank[i] = span
		} else {
			rank[i] = rank[i-1]
		}
	}
	return backList, rank
}

func (s *Skiplist) del(node *SkiplistNode, backList []*SkiplistNode) {
	if backList == nil {
		nbackList, _ := s.getBackList(node)
		backList = nbackList
	}

	for i := 0; i < len(node.levels); i++ {
		back := backList[i]
		back.levels[i].span += node.levels[i].span - 1
		next := node.levels[i].forward
		back.levels[i].forward = next
		if i == 0 {
			if next != nil {
				next.backward = back
				back.levels[i].span = 1
			} else {
				s.tail = back
				back.levels[i].span = 0
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
	backList, rank := s.getBackList(node)
	if h > s.level {
		for i := s.level; i < h; i++ {
			backList = append(backList, s.head)
		}
		s.level = h
	}

	if backList[0].levels[0].forward != nil &&
		backList[0].levels[0].forward.score == score &&
		backList[0].levels[0].forward.ele == key {
		s.del(node, backList)
	}

	for i := range h {
		node.levels[i].forward = backList[i].levels[i].forward
		backList[i].levels[i].forward = node
		backList[i].levels[i].span = rank[i] + 1
	}

	node.backward = backList[0]
	if node.levels[0].forward != nil {
		node.levels[0].forward.backward = node
		node.levels[0].span = 1
	} else {
		s.tail = node
		node.levels[0].span = 1
	}
	s.length++
	return nil, false
}
