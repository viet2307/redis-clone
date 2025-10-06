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
	span := uint32(0)
	for i := s.level - 1; i >= 0; i-- {
		next := curr.levels[i].forward
		for next != nil && (next.score < node.score || (next.score == node.score && next.ele < node.ele)) {
			span += curr.levels[i].span
			curr = next
			next = curr.levels[i].forward
		}
		backList[i] = curr
		rank[i] = span
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
		next := node.levels[i].forward
		back.levels[i].forward = next

		back.levels[i].span += node.levels[i].span - 1
		if i == 0 {
			if next != nil {
				next.backward = back
				continue
			}
			s.tail = back
			if back == s.tail {
				s.tail = nil
			}
		}
	}

	for i := len(node.levels); i < s.level; i++ {
		backList[i].levels[i].span--
	}

	for s.level > 1 && s.head.levels[s.level-1].forward == nil {
		s.level--
	}

	s.length--
}

func (s *Skiplist) Zadd(key string, score float64) (interface{}, bool) {
	h := s.coinFlip()
	node := s.newNode(key, score, h)
	backList, rank := s.getBackList(node)
	if h > s.level {
		for i := s.level; i < h; i++ {
			backList = append(backList, s.head)
			rank = append(rank, rank[s.level-1])
		}
		s.level = h
	}

	if backList[0].levels[0].forward != nil &&
		backList[0].levels[0].forward.score == score &&
		backList[0].levels[0].forward.ele == key {
		oldNode := backList[0].levels[0].forward
		s.del(oldNode, backList)
	}

	for i := range h {
		next := backList[i].levels[i].forward
		node.levels[i].forward = next
		backList[i].levels[i].forward = node

		switch i {
		case 0:
			backList[i].levels[i].span = 1
			node.backward = backList[i]
			if next != nil {
				node.levels[i].span = 1
				next.backward = node
				continue
			}

			node.levels[i].span = 0
			s.tail = node
		default:
			oldSpan := backList[i].levels[i].span
			backList[i].levels[i].span = rank[i] - rank[0] + 1
			if next != nil {
				node.levels[i].span = oldSpan - (rank[i] - rank[0])
				continue
			}

			node.levels[i].span = 0
		}
	}

	for i := h; i < s.level; i++ {
		backList[i].levels[i].span++
	}
	s.length++
	return nil, false
}
