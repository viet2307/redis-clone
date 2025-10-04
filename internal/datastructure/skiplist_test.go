package datastructure

import (
	"math/rand"
	"testing"
)

// Helper to seed randomness for deterministic tests
func seedRand() {
	rand.Seed(42)
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

func TestCoinFlip(t *testing.T) {
	seedRand()
	sl := NewSkiplist(16)

	// Should always return at least 1
	h := sl.coinFlip()
	if h < 1 || h > sl.maxLevel {
		t.Errorf("coinFlip returned invalid height %d", h)
	}
}

func TestNewNode(t *testing.T) {
	sl := NewSkiplist(8)
	node := sl.newNode("foo", 3.14, 4)

	if node.ele != "foo" {
		t.Errorf("Expected ele 'foo', got %s", node.ele)
	}
	if node.score != 3.14 {
		t.Errorf("Expected score 3.14, got %f", node.score)
	}
	if len(node.levels) != 4 {
		t.Errorf("Expected 4 levels, got %d", len(node.levels))
	}
}

func TestGetBackList(t *testing.T) {
	sl := NewSkiplist(4)
	sl.level = 2
	_, _ = sl.Zadd("a", 1.0)
	_, _ = sl.Zadd("b", 2.0)
	_, _ = sl.Zadd("c", 3.0)
	node := sl.newNode("d", 2.5, 2)
	backList, _ := sl.getBackList(node)

	if len(backList) != sl.level {
		t.Errorf("BackList size mismatch: expected %d, got %d", sl.level, len(backList))
	}
	for i, back := range backList {
		if back == nil {
			t.Errorf("BackList[%d] should not be nil", i)
		}
	}
}

func TestZaddAndDel(t *testing.T) {
	sl := NewSkiplist(6)
	_, _ = sl.Zadd("x", 10.0)
	_, _ = sl.Zadd("y", 20.0)
	_, _ = sl.Zadd("z", 30.0)

	// Check length
	if sl.length != 3 {
		t.Errorf("Expected length 3, got %d", sl.length)
	}

	// Check tail
	if sl.tail == nil || sl.tail.ele != "z" {
		t.Errorf("Expected tail to be 'z', got '%v'", sl.tail)
	}

	_, _ = sl.Zadd("y", 20.0)

	// Remove tail node manually
	node := sl.tail
	backList, _ := sl.getBackList(node)
	sl.del(node, backList)
	if sl.length != 2 {
		t.Errorf("After deletion, expected length 1, got %d", sl.length)
	}
}

func TestBackwardLinks(t *testing.T) {
	sl := NewSkiplist(4)
	_, _ = sl.Zadd("A", 1.0)
	_, _ = sl.Zadd("B", 2.0)
	_, _ = sl.Zadd("C", 3.0)

	nodeB := sl.head.levels[0].forward.levels[0].forward // Should be "B"
	if nodeB.backward == nil || nodeB.backward.ele != "A" {
		t.Errorf("Backward link from B should be A, got %v", nodeB.backward)
	}
	nodeC := sl.tail
	if nodeC.backward == nil || nodeC.backward.ele != "B" {
		t.Errorf("Backward link from C should be B, got %v", nodeC.backward)
	}
}
