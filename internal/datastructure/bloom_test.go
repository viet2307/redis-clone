package datastructure

import (
	"fmt"
	"testing"
)

func TestEmptyBloom(t *testing.T) {
	bloom := NewBloom(0.01, 1000)
	tests := []string{"test1", "test2", "test3"}

	for _, item := range tests {
		if bloom.Exist(item) {
			t.Errorf("Empty bloom filter returned true for %q", item)
		}
	}
}

func TestBloomAdd(t *testing.T) {
	bloom := NewBloom(0.01, 1000)
	bloom.Add("banana")
	bloom.Add("banana")
	bloom.Add("banana2")

	tests := []struct {
		key   string
		exist bool
	}{
		{key: "banana", exist: true},
		{key: "banana", exist: true},
		{key: "banana", exist: true},
		{key: "banana2", exist: true},
		{key: "banana3", exist: false},
		{key: "bababana", exist: false},
		{key: "babanana", exist: false},
	}

	for _, tt := range tests {
		ok := bloom.Exist(tt.key)
		if ok != tt.exist && !ok {
			t.Errorf("ERROR, bloom filter return false negative for key %s: want %v, got %v", tt.key, tt.exist, ok)
		}
	}
}

func TestLargeScale(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large scale test in short mode")
	}

	entries := uint64(100000)
	bloom := NewBloom(0.001, entries)

	for i := 0; i < int(entries); i++ {
		bloom.Add(fmt.Sprintf("Item number %d", i))
	}

	sampleSize := 1000
	for i := range sampleSize {
		idx := i * (int(entries) / sampleSize)
		if !bloom.Exist(fmt.Sprintf("Item number %d", idx)) {
			t.Errorf("Lost item at index %d", idx)
		}
	}

	t.Logf("Success")
}

func BenchmarkAdd(b *testing.B) {
	bloom := NewBloom(0.01, uint64(b.N))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bloom.Add(fmt.Sprintf("item number %d", i))
	}
}

func BenchmarkExist(b *testing.B) {
	bloom := NewBloom(0.01, 10000)

	for i := range 10000 {
		bloom.Add(fmt.Sprintf("item number %d", i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bloom.Exist(fmt.Sprintf("item_%d", i%10000))
	}
}
