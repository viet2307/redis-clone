package datastructure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

const Log10PointFive = -0.30102999566

type CMS struct {
	w, d    uint32
	counter [][]uint32
}

/*
Get w, d based on errRate and errProb
errRate: upper bound of overcount (per element)
errProb: probability of error count happening
*/
func CalcCMSDim(errRate float64, errProb float64) (uint32, uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(errProb) / Log10PointFive))
	return w, d
}

func (c *CMS) hashfunc(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}

func NewCMS(errRate float64, errProb float64) *CMS {
	w, d := CalcCMSDim(errRate, errProb)
	counter := make([][]uint32, d)
	for i := uint32(0); i < d; i++ {
		counter[i] = make([]uint32, w)
	}
	return &CMS{
		w:       w,
		d:       d,
		counter: counter,
	}
}

func (c *CMS) IncrBy(item string, value uint32) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := uint32(0); i < c.d; i++ {
		hash := c.hashfunc(item, i)
		j := hash % c.w

		if math.MaxUint32-c.counter[i][j] < value {
			c.counter[i][j] = math.MaxUint32
		} else {
			c.counter[i][j] += value
		}

		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}
	return minCount
}

func (c *CMS) Query(item string) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := uint32(0); i < c.d; i++ {
		hash := c.hashfunc(item, i)

		j := hash % c.w

		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}
	return minCount
}
