package lsm

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SSTable struct {
	MinKey string            `json:"min_key"`
	MaxKey string            `json:"max_key"`
	Data   map[string]string `json:"data"`
	Size   int               `json:"size"`
	Level  int               `json:"level"` // NEW: Added Level field
}

func NewSSTable(data map[string]string, level int) *SSTable {
	minKey, maxKey := findMinMaxKeys(data)
	return &SSTable{
		MinKey: minKey,
		MaxKey: maxKey,
		Data:   data,
		Size:   len(data),
		Level:  level,
	}
}

func findMinMaxKeys(data map[string]string) (string, string) {
	minKey := ""
	maxKey := ""
	for k := range data {
		if minKey == "" || k < minKey {
			minKey = k
		}
		if k > maxKey {
			maxKey = k
		}
	}
	return minKey, maxKey
}

func (sst *SSTable) SaveToDisk(filename string) error {
	path := filepath.Clean(filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(sst)
}

// load sstable from file
func LoadSSTable(path string) (*SSTable, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sst SSTable
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sst)
	return &sst, err
}

func (sst *SSTable) Get(key string) (string, bool) {
	val, found := sst.Data[key]
	return val, found
}

func (sst *SSTable) Metadata() map[string]interface{} {
	var minKey, maxKey string
	for k := range sst.Data {
		if minKey == "" || k < minKey {
			minKey = k
		}
		if maxKey == "" || k > maxKey {
			maxKey = k
		}
	}
	return map[string]interface{}{
		"minKey": sst.MinKey,
		"maxKey": sst.MaxKey,
		"size":   len(sst.Data),
		"level":  sst.Level,
	}
}
