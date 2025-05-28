package lsm

import (
	"errors"
	"sort"
	"sync"
)

type Memtable struct {
	mu    sync.RWMutex
	table map[string]string
	keys  []string // sorted keys for iteration
}

func NewMemtable() *Memtable {
	return &Memtable{
		table: make(map[string]string),
		keys:  []string{},
	}
}

// Put adds or updates a key-value pair
func (m *Memtable) Put(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.table[key]; !exists {
		m.keys = append(m.keys, key)
		sort.Strings(m.keys)
	}
	m.table[key] = value
}

// Get retrieves the value for a key
func (m *Memtable) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.table[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

// Flush returns all entries and resets the memtable
func (m *Memtable) Flush() map[string]string {
	m.mu.Lock()
	defer m.mu.Unlock()

	flushed := m.table
	m.table = make(map[string]string)
	m.keys = []string{}
	return flushed
}
