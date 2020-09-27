package bloom

import (
	"hash"
	"hash/fnv"

	"github.com/spaolacci/murmur3"
)

// BloomFilter probabilistic data structure definition
type BloomFilter struct {
	// Please use math/big as an alternative for bitset
	bitset  []int         // The bloom-filter bitset
	n       uint          // Number of elements in the filter
	m       uint          // Size of the bloom filter
	hashfns []hash.Hash64 // The hash functions
}

// DefaultHashFunctions for BloomFilter
var DefaultHashFunctions = []hash.Hash64{murmur3.New64(), fnv.New64(), fnv.New64a()}

// Returns a new BloomFilter object,
func NewBloomFilter(m uint, hashes []hash.Hash64) *BloomFilter {
	if hashes == nil || len(hashes) <= 0 {
		return &BloomFilter{
			bitset:  make([]int, m),
			m:       m,
			n:       uint(0),
			hashfns: DefaultHashFunctions,
		}
	}
	return &BloomFilter{
		bitset:  make([]int, m),
		m:       m,
		n:       uint(0),
		hashfns: hashes,
	}
}

// Adds the item into the bloom filter set by hashing in over the hash functions
func (bf *BloomFilter) Add(item []byte) {
	for _, v := range bf.convertToHashes(item) {
		// hash result example = "131313123123131233123123" want to fit in the vector bit? Use modulo
		position := uint(v) % bf.m
		// position must be between 0 and (m - 1)
		bf.bitset[position] = 1
	}
	// ex: hash results after module: 0, 2, 3 ->> [1, 0, 1, 1]
	bf.n += 1
}

// Test if the item into the bloom filter is set by hashing in over the hash functions
func (bf *BloomFilter) Test(item []byte) bool {
	for _, v := range bf.convertToHashes(item) {
		position := uint(v) % bf.m
		if bf.bitset[position] != 1 {
			return false
		}
	}
	return true
}

// Calculates all the hash values by hashing in over the hash functions
func (bf *BloomFilter) convertToHashes(item []byte) []uint64 {
	// result length depends on total hash functions
	var result []uint64

	for _, hashFunc := range bf.hashfns {
		hashFunc.Write(item)
		result = append(result, hashFunc.Sum64())
		hashFunc.Reset()
	}

	return result
}
