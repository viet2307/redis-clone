package datastructure

import (
	"math/rand"
	"strings"
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

func NewZset() *ZSet {
	return &ZSet{
		zskiplist: NewSkiplist(32),
		dict:      make(map[string]float64),
	}
}

func NewSkiplist(maxLevel int) *Skiplist {
	skiplistLevel := make([]SkiplistLevel, maxLevel)
	for i := range skiplistLevel {
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

func newNode(ele string, score float64, h int) *SkiplistNode {
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

func (s *Skiplist) skiplistDel(node *SkiplistNode, backList []*SkiplistNode) {
	if backList == nil {
		backList, _ = s.getBackList(node)
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
			if back == s.head {
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

func (z *ZSet) zsetDel(node *SkiplistNode, backList []*SkiplistNode) {
	z.zskiplist.skiplistDel(node, backList)
	delete(z.dict, node.ele)
}

func (s *Skiplist) skiplistAdd(ele string, score float64) {
	h := s.coinFlip()
	node := newNode(ele, score, h)
	backList, rank := s.getBackList(node)
	if h > s.level {
		for i := s.level; i < h; i++ {
			backList = append(backList, s.head)
			rank = append(rank, rank[s.level-1])
		}
		s.level = h
	}

	node.backward = backList[0]

	for i := 0; i < h; i++ {
		next := backList[i].levels[i].forward
		node.levels[i].forward = next
		backList[i].levels[i].forward = node

		if next != nil {
			if i == 0 {
				backList[i].levels[i].span = 1
				node.levels[i].span = 1
			} else {
				oldSpan := backList[i].levels[i].span
				backList[i].levels[i].span = rank[i] - rank[0] + 1
				node.levels[i].span = oldSpan - (rank[i] - rank[0])
			}
			next.backward = node
		} else {
			if i == 0 {
				backList[i].levels[i].span = 1
			} else {
				backList[i].levels[i].span = rank[i] - rank[0] + 1
			}
			node.levels[i].span = 0
			s.tail = node
		}
	}

	for i := h; i < s.level; i++ {
		backList[i].levels[i].span++
	}
	s.length++
}

func (z *ZSet) Zadd(ele string, score float64) (int, bool) {
	if oldScore, exists := z.dict[ele]; exists {
		if oldScore == score {
			return 0, false
		}
		oldNode := newNode(ele, oldScore, 1)
		backList, _ := z.zskiplist.getBackList(oldNode)
		if backList[0].levels[0].forward != nil &&
			backList[0].levels[0].forward.ele == oldNode.ele {
			z.zsetDel(backList[0].levels[0].forward, backList)
		}
	}

	z.dict[ele] = score
	z.zskiplist.skiplistAdd(ele, score)
	return 1, false
}

func (z *ZSet) Zscore(ele string) (interface{}, bool) {
	score, exists := z.dict[ele]
	if !exists {
		return nil, false
	}
	return float64(score), false
}

func (z *ZSet) Zrank(ele string) (int, bool) {
	s := z.zskiplist
	score, exists := z.dict[ele]
	if !exists {
		return -1, false
	}

	rank := uint32(0)
	curr := s.head
	for l := s.level - 1; l >= 0; l-- {

		for next := curr.levels[l].forward; next != nil && (next.score < score ||
			(next.score == score &&
				strings.Compare(next.ele, ele) < 0)); {
			rank += curr.levels[l].span
			curr = next
			next = curr.levels[l].forward
		}
	}
	next := curr.levels[0].forward
	if next != nil && next.score == score && strings.Compare(next.ele, ele) == 0 {
		return int(rank), false
	}

	return -1, false
}
