package lsm

import (
	"fmt"
	"sort"
)

func (tree *LSMTree) compactLevel(level int) {
	tree.compactionLogs = append(tree.compactionLogs, fmt.Sprintf("Starting compaction at level %d", level))

	if level+1 >= len(tree.sstables) {
		tree.sstables = append(tree.sstables, []*SSTable{})
	}

	// check if compaction is needed
	if len(tree.sstables[level]) < tree.levelMaxFiles {
		return
	}

	toCompact := tree.sstables[level]
	tree.sstables[level] = nil

	mergedData := make(map[string]string)
	for _, sst := range toCompact {
		for k, v := range sst.Data {
			mergedData[k] = v
		}
	}

	// sort keys for compaction
	var keys []string
	for k := range mergedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// create new SStable for level+1
	newSST := NewSSTable(mergedData, 0)
	tree.sstables[level+1] = append(tree.sstables[level+1], newSST)

	tree.compactionLogs = append(tree.compactionLogs, fmt.Sprintf("Compacted %d SSTables from level %d into 1 SSTable at level %d with %d keys", len(toCompact), level, level+1, len(mergedData)))

	// recursively compact next level if needed
	tree.compactLevel(level + 1)
}
