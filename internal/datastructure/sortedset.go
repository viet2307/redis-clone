package datastructure

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
	head   *SkiplistNode
	tail   *SkiplistNode
	length uint32
	level  int
}

type ZSet struct {
	zskiplist *Skiplist
	dict      map[string]float64
}
