package lsm

type BloomFilter struct {
	bits []bool
	size int
}

func NewBloomFilter(size int) *BloomFilter {
	return &BloomFilter{
		bits: make([]bool, size),
		size: size,
	}
}

// Very simple hash function for demonstration (not production quality)
func hash(key string, seed int) int {
	h := 0
	for i := 0; i < len(key); i++ {
		h = (h*31 + int(key[i]) + seed) % 100000
	}
	return h
}

func (bf *BloomFilter) Add(key string) {
	idx := hash(key, 1) % bf.size
	bf.bits[idx] = true
}

func (bf *BloomFilter) MightContain(key string) bool {
	idx := hash(key, 1) % bf.size
	return bf.bits[idx]
}
