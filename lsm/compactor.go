package lsm

import (
	"fmt"
)

// Merge two SSTables (simulate compaction)
// New SSTable path must be provided
func CompactSSTables(sst1, sst2 *SSTable, newPath string) (*SSTable, error) {
	mergedData := make(map[string]string)

	// Add all keys from sst1
	for k, v := range sst1.Data {
		mergedData[k] = v
	}
	// Add/overwrite keys from sst2 (newer data overwrites older)
	for k, v := range sst2.Data {
		mergedData[k] = v
	}

	fmt.Printf("Compacting SSTables into %s (total keys: %d)\n", newPath, len(mergedData))

	return NewSSTableFromMap(newPath, mergedData)
}
