package lsm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LSMTree struct {
	memtable       *Memtable
	memtableCap    int
	bloom          *BloomFilter
	sstables       [][]*SSTable
	levelMaxFiles  int
	compactionLogs []string
}

func NewLSMTree(memtableCap, bloomSize int) *LSMTree {
	return &LSMTree{
		memtable:      NewMemtable(memtableCap),
		memtableCap:   memtableCap,
		bloom:         NewBloomFilter(uint(bloomSize)),
		sstables:      [][]*SSTable{{}},
		levelMaxFiles: 4,
	}
}

func (tree *LSMTree) Put(key, val string) {
	tree.memtable.Put(key, val)
	tree.bloom.Add(key)

	if tree.memtable.Size() >= tree.memtableCap {
		tree.flushMemtable()
	}
}

func (tree *LSMTree) flushMemtable() {
	data := tree.memtable.Flush()
	newSST := NewSSTable(data, 0)
	tree.sstables[0] = append(tree.sstables[0], newSST)
	tree.compactionLogs = append(tree.compactionLogs, "Memtable flushed to SSTable level 0 with "+fmt.Sprint(len(data))+" keys")
	tree.compactLevel(0)

	filename := fmt.Sprintf("sstables/level%d-%d.sst", newSST.Level, time.Now().UnixNano())
	err := newSST.SaveToDisk(filename)
	if err != nil {
		log.Println("Failed to persist SSTable:", err)
	}
}

func (tree *LSMTree) Get(key string) (string, bool) {
	// Check memtable
	if val, ok := tree.memtable.Get(key); ok {
		return val, true
	}
	// Check bloom filter first
	if !tree.bloom.Test(key) {
		return "", false
	}
	// Check SSTables from level 0 up
	for _, level := range tree.sstables {
		for _, sst := range level {
			if val, ok := sst.Get(key); ok {
				return val, true
			}
		}
	}
	return "", false
}

func (tree *LSMTree) MemtableKeys() []string {
	return tree.memtable.Keys()
}

func (tree *LSMTree) BloomStats() map[string]interface{} {
	return tree.bloom.Stats()
}

func (tree *LSMTree) SSTablesMetadata() []map[string]interface{} {
	levels := []map[string]interface{}{}
	for i, level := range tree.sstables {
		var levelInfo []map[string]interface{}
		for _, sst := range level {
			levelInfo = append(levelInfo, sst.Metadata())
		}
		levels = append(levels, map[string]interface{}{
			"level": i,
			"ssts":  levelInfo,
		})
	}
	return levels
}

func (tree *LSMTree) CompactionLogs() []string {
	return tree.compactionLogs
}

func (lsm *LSMTree) LoadSSTablesFromDisk() {
	files, _ := filepath.Glob("sstables/*.sst")
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var sst SSTable
		if err := json.Unmarshal(data, &sst); err == nil {
			lsm.sstables[sst.Level] = append(lsm.sstables[sst.Level], &sst)
		}
	}
}
