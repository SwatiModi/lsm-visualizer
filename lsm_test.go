package main

import (
	"fmt"
	"testing"
)

// TestPutAndGet tests basic Put and Get operations in the Memtable and SSTables
func TestPutAndGet(t *testing.T) {
	db := NewLSMTree()

	// Put key-value pairs
	keys := []string{"key1", "key2", "key3"}
	values := []string{"val1", "val2", "val3"}

	for i, k := range keys {
		err := db.Put(k, values[i])
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
	}

	// Get values for existing keys
	for i, k := range keys {
		v, err := db.Get(k)
		if err != nil {
			t.Fatalf("Get failed for key %s: %v", k, err)
		}
		if v != values[i] {
			t.Errorf("Expected %s for key %s, got %s", values[i], k, v)
		}
	}

	// Get for non-existing key
	_, err := db.Get("nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existing key")
	}
}

// TestFlushMemtable tests that the memtable flushes to SSTable correctly
func TestFlushMemtable(t *testing.T) {
	db := NewLSMTree()
	// Put enough keys to trigger flush (threshold is 5)
	for i := 0; i < 5; i++ {
		key := "key" + fmt.Sprint('a'+i)
		err := db.Put(key, "val")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
	}

	if len(db.memtable) != 0 {
		t.Errorf("Memtable should be empty after flush, got size %d", len(db.memtable))
	}
	if len(db.sstables) == 0 {
		t.Errorf("Expected at least 1 SSTable after flush")
	}
}

// TestCompaction tests that compaction merges SSTables correctly
func TestCompaction(t *testing.T) {
	db := NewLSMTree()

	// Add enough entries to create 2 SSTables
	for i := 0; i < 10; i++ {
		key := "key" + fmt.Sprint('a'+(i%26))
		err := db.Put(key, "val")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
	}

	if len(db.sstables) < 2 {
		t.Fatalf("Expected at least 2 SSTables before compaction, got %d", len(db.sstables))
	}

	err := db.Compact()
	if err != nil {
		t.Fatalf("Compaction failed: %v", err)
	}

	if len(db.sstables) != 1 {
		t.Errorf("Expected 1 SSTable after compaction, got %d", len(db.sstables))
	}
}
