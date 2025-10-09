package datastructure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

const SEED uint32 = 19980605

var LogSquare float64 = math.Log(2) * math.Log(2)

type Bloom struct {
	hashes        int
	errorRate     float64
	entries       uint64
	bitPerEntries float64
	bf            []byte
	bits          uint64
	bytes         uint64
}

type HashVal struct {
	a, b uint64
}

func calBitsPerEntries(errorRate float64) float64 {
	num := math.Log(errorRate)
	return math.Abs(-(num / LogSquare))
}

func NewBloom(errorRate float64, entries uint64) *Bloom {
	bloom := Bloom{
		entries:   entries,
		errorRate: errorRate,
	}
	bloom.bitPerEntries = calBitsPerEntries(errorRate)
	bits := entries * uint64(bloom.bitPerEntries)
	if bits%64 != 0 {
		bloom.bytes = ((bits / 64) + 1) * 8
	} else {
		bloom.bytes = bits / 8
	}
	bloom.bits = bloom.bytes * 8
	bloom.hashes = int(math.Ceil(math.Log(2) * bloom.bitPerEntries))
	bloom.bf = make([]byte, bloom.bytes)
	return &bloom
}

func calHash(entry string) HashVal {
	hasher := murmur3.New128WithSeed(SEED)
	hasher.Write([]byte(entry))
	a, b := hasher.Sum128()
	return HashVal{
		a: a,
		b: b,
	}
}

func (b *Bloom) Add(entry string) {
	initHash := calHash(entry)
	for i := range b.hashes {
		hash := (initHash.a + initHash.b*uint64(i)) % 8
		b.bf[hash/8] |= 1 << (hash / 8)
	}
}

func (b *Bloom) Exist(entry string) bool {
	initHash := calHash(entry)
	for i := range b.hashes {
		hash := (initHash.a + initHash.b*uint64(i)) % 8
		if b.bf[hash/8]&1<<(hash/8) == 0 {
			return false
		}
	}
	return true
}
