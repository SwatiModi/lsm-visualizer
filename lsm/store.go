package lsm

import (
	"fmt"
	"os"
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	memtable *Memtable
	wal      *WAL

	sstables []*SSTable // level 0 for now

	bloom *BloomFilter

	sstCounter int
}

func NewStore(walPath string) (*Store, error) {
	wal, err := NewWAL(walPath)
	if err != nil {
		return nil, err
	}
	return &Store{
		memtable: NewMemtable(),
		wal:      wal,
		bloom:    NewBloomFilter(1000),
		sstables: []*SSTable{},
	}, nil
}

func (s *Store) Put(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.memtable.Put(key, value)
	err := s.wal.Write(key, value)
	if err != nil {
		return err
	}
	s.bloom.Add(key)
	// Flush if memtable size exceeds threshold (e.g., 5 entries for demo)
	if len(s.memtable.keys) >= 5 {
		return s.flushMemtable()
	}
	return nil
}

func (s *Store) flushMemtable() error {
	flushed := s.memtable.Flush()
	sstPath := fmt.Sprintf("sstable_%d.json", s.sstCounter)
	s.sstCounter++
	sst, err := NewSSTableFromMap(sstPath, flushed)
	if err != nil {
		return err
	}
	s.sstables = append(s.sstables, sst)
	fmt.Println("Flushed Memtable to SSTable:", sstPath)
	return nil
}

func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check Memtable first
	val, err := s.memtable.Get(key)
	if err == nil {
		return val, nil
	}

	// Check Bloom filter
	if !s.bloom.MightContain(key) {
		return "", fmt.Errorf("key %s definitely not present (bloom filter miss)", key)
	}

	// Search SSTables in reverse order (newest first)
	for i := len(s.sstables) - 1; i >= 0; i-- {
		val, err := s.sstables[i].Get(key)
		if err == nil {
			return val, nil
		}
	}
	return "", fmt.Errorf("key %s not found", key)
}

func (s *Store) Close() error {
	return s.wal.Close()
}

// Compact SSTables example (compact last 2 SSTables)
func (s *Store) Compact() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sstables) < 2 {
		return nil
	}
	sst1 := s.sstables[len(s.sstables)-2]
	sst2 := s.sstables[len(s.sstables)-1]

	newPath := fmt.Sprintf("sstable_%d_compacted.json", s.sstCounter)
	s.sstCounter++

	newSST, err := CompactSSTables(sst1, sst2, newPath)
	if err != nil {
		return err
	}

	// Remove last two SSTables and add the new one
	s.sstables = s.sstables[:len(s.sstables)-2]
	s.sstables = append(s.sstables, newSST)

	// Optionally remove old SSTable files from disk
	os.Remove(sst1.Path)
	os.Remove(sst2.Path)

	fmt.Println("Compacted SSTables into:", newPath)
	return nil
}
