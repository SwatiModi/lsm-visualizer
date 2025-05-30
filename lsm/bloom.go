package lsm

import "github.com/willf/bloom"

type BloomFilter struct {
	filter *bloom.BloomFilter
}

func NewBloomFilter(size uint) *BloomFilter {
	return &BloomFilter{
		filter: bloom.NewWithEstimates(size, 0.01),
	}
}

func (bf *BloomFilter) Add(key string) {
	bf.filter.AddString(key)
}

func (bf *BloomFilter) Test(key string) bool {
	return bf.filter.TestString(key)
}

func (bf *BloomFilter) Stats() map[string]interface{} {
	return map[string]interface{}{
		"capacity": bf.filter.Cap(),
		"k":        bf.filter.K(),
		"est_fp":   bf.filter.EstimateFalsePositiveRate(5),
	}
}
