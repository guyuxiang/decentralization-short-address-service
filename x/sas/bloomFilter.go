package sas

import (
	"github.com/spaolacci/murmur3"
	"math/rand"
)

var GlobalBloomFilter *BloomFilter

type BitMap struct {
	bits []byte
	vmax uint
}

func NewBitMap(maxVal uint) *BitMap {
	var max uint = 8192
	if maxVal > 0 && maxVal > max {
		max = maxVal
	}
	sz := (max + 7) / 8
	return &BitMap{
		bits: make([]byte, sz, sz),
		vmax: max,
	}
}

func (bm *BitMap) Set(num uint) {
	if num > bm.vmax {
		bm.vmax += 1024
		if bm.vmax < num {
			bm.vmax = num
		}
		dd := int(num+7)/8 - len(bm.bits)
		if dd > 0 {
			bm.bits = append(bm.bits, make([]byte, dd, dd)...)
		}
	}
	bm.bits[num/8] = bm.bits[num/8] | (1 << (num % 8))
}

func (bm *BitMap) Check(num uint) bool {
	if num > bm.vmax {
		return false
	}
	return bm.bits[num/8]&(1<<(num%8)) != 0
}

type BloomFilter struct {
	bset  *BitMap
	size  uint
	seeds []uint32
}

func NewBloomFilter(sizeVal uint, seedNum int) *BloomFilter {
	var size uint = 1024 * 1024
	if sizeVal > 0 && sizeVal > size {
		size = sizeVal
	}

	bf := &BloomFilter{
		bset: NewBitMap(size),
		size: size,
	}

	for i := 0; i < seedNum; i++ {
		bf.seeds = append(bf.seeds, uint32(rand.Intn(100)))
	}

	return bf
}

func (bf *BloomFilter) Check(value string) bool {
	for _, seed := range bf.seeds {
		hash := murmur3.Sum64WithSeed([]byte(value), seed)
		hash = hash % uint64(bf.size)
		ret := bf.bset.Check(uint(hash))
		if !ret {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) Set(value string) {
	for _, seed := range bf.seeds {
		hash := murmur3.Sum64WithSeed([]byte(value), seed)
		hash = hash % uint64(bf.size)
		bf.bset.Set(uint(hash))
	}
}
