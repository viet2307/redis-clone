package datastructure

import (
	"fmt"
	"math/rand"
	"testing"
)

func seedRand() {
	rand.Seed(42)
}

func TestZrankEmptyZSet(t *testing.T) {
	z := NewZset()
	rank := z.Zrank("nonexistent")
	if rank != -1 {
		t.Errorf("Expected rank -1 for nonexistent element in empty zset, got %v", rank)
	}
}

func TestNewSkiplist(t *testing.T) {
	maxLevel := 6
	sl := NewSkiplist(maxLevel)

	if sl.maxLevel != maxLevel {
		t.Errorf("Expected maxLevel %d, got %d", maxLevel, sl.maxLevel)
	}
	if sl.head == nil {
		t.Error("Head node should not be nil")
	}
	if len(sl.head.levels) != maxLevel {
		t.Errorf("Head node should have %d levels, got %d", maxLevel, len(sl.head.levels))
	}
	if sl.level != 1 {
		t.Errorf("Initial skiplist level should be 1, got %d", sl.level)
	}
}

func TestZrankDuplicateScores(t *testing.T) {
	z := NewZset()
	z.Zadd("zebra", 100.0)
	z.Zadd("apple", 100.0)
	z.Zadd("monkey", 100.0)
	z.Zadd("banana", 100.0)

	tests := []struct {
		ele  string
		want int
	}{
		{"apple", 0},
		{"banana", 1},
		{"monkey", 2},
		{"zebra", 3},
	}

	for _, tt := range tests {
		rank := z.Zrank(tt.ele)
		if rank != tt.want {
			t.Errorf("Zrank(%s) = %v, want %v", tt.ele, rank, tt.want)
		}
		fmt.Printf("After update: Zrank(%s) = %v, want %v\n", tt.ele, rank, tt.want)
	}
}

func TestZrankAfterScoreUpdate(t *testing.T) {
	z := NewZset()
	z.Zadd("alice", 100.0)
	z.Zadd("bob", 200.0)
	z.Zadd("charlie", 150.0)

	z.Zadd("bob", 50.0)

	tests := []struct {
		ele  string
		want int
	}{
		{"bob", 0},
		{"alice", 1},
		{"charlie", 2},
	}

	for _, tt := range tests {
		rank := z.Zrank(tt.ele)
		if rank != tt.want {
			t.Errorf("After update: Zrank(%s) = %v, want %v", tt.ele, rank, tt.want)
		}
		fmt.Printf("After update: Zrank(%s) = %v, want %v\n", tt.ele, rank, tt.want)
	}
}

func TestCoinFlip(t *testing.T) {
	seedRand()
	sl := NewSkiplist(16)

	h := sl.coinFlip()
	if h < 1 || h > sl.maxLevel {
		t.Errorf("coinFlip returned invalid height %d", h)
	}
}

func TestZaddAndDel(t *testing.T) {
	z := NewZset()
	s := z.zskiplist

	z.Zadd("x", 10.0)
	z.Zadd("y", 20.0)
	z.Zadd("z", 30.0)

	if s.length != 3 {
		t.Errorf("Expected length 3, got %d", s.length)
	}

	if s.tail == nil || s.tail.ele != "z" {
		t.Errorf("Expected tail to be 'z', got '%v'", s.tail)
	}

	z.Zadd("y", 20.0)
	if s.length != 3 {
		t.Errorf("After update, expected length 3, got %d", s.length)
	}

	node := s.tail
	backList, _ := s.getBackList(node)
	z.zskiplist.skiplistDel(node, backList)
	if s.length != 2 {
		t.Errorf("After deletion, expected length 2, got %d", s.length)
	}
}

func TestBackwardLinks(t *testing.T) {
	z := NewZset()
	s := z.zskiplist

	z.Zadd("A", 1.0)
	z.Zadd("B", 2.0)
	z.Zadd("C", 3.0)

	nodeB := s.head.levels[0].forward.levels[0].forward
	if nodeB.backward == nil || nodeB.backward.ele != "A" {
		t.Errorf("Backward link from B should be A, got %v", nodeB.backward)
	}
	nodeC := s.tail
	if nodeC.backward == nil || nodeC.backward.ele != "B" {
		t.Errorf("Backward link from C should be B, got %v", nodeC.backward)
	}
}

func TestZscore(t *testing.T) {
	z := NewZset()
	z.Zadd("alice", 100.0)
	z.Zadd("bob", 200.0)

	score := z.Zscore("alice")
	if score == nil || score != 100.0 {
		t.Errorf("Expected score 100.0 for alice, got %v", score)
	}

	score = z.Zscore("charlie")
	if score != nil {
		t.Error("Expected nil for non-existent member")
	}

	z.Zadd("alice", 100.5)
	z.Zadd("bob", 200.7)

	score = z.Zscore("alice")
	if score == nil || score != 100.5 {
		t.Errorf("Expected score 100.5 for alice, got %v", score)
	}

	score = z.Zscore("bob")
	if score == nil || score != 200.7 {
		t.Errorf("Expected score 200.7 for alice, got %v", score)
	}

	score = z.Zscore("charlie")
	if score != nil {
		t.Error("Expected nil for non-existent member")
	}
}

func TestZrank(t *testing.T) {
	z := NewZset()
	z.Zadd("alice", 100.0)
	z.Zadd("bob", 200.0)
	z.Zadd("charlie", 150.0)

	rank := z.Zrank("alice")
	if rank == -1 || rank != 0 {
		t.Errorf("Expected rank 0 for alice, got %v", rank)
	}

	rank = z.Zrank("charlie")
	if rank == -1 || rank != 1 {
		t.Errorf("Expected rank 1 for charlie, got %v", rank)
	}

	rank = z.Zrank("bob")
	if rank == -1 || rank != 2 {
		t.Errorf("Expected rank 2 for bob, got %v", rank)
	}
}
