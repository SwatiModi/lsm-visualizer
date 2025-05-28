package lsm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"sync"
)

type SSTable struct {
	mu   sync.RWMutex
	Path string
	Data map[string]string
	Keys []string
}

func NewSSTableFromMap(path string, data map[string]string) (*SSTable, error) {
	sst := &SSTable{
		Path: path,
		Data: data,
		Keys: make([]string, 0, len(data)),
	}
	for k := range data {
		sst.Keys = append(sst.Keys, k)
	}
	sort.Strings(sst.Keys)

	// Save to disk as JSON
	err := sst.save()
	if err != nil {
		return nil, err
	}

	return sst, nil
}

func (s *SSTable) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonData, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, jsonData, 0644)
}

func LoadSSTable(path string) (*SSTable, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data map[string]string
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	sst := &SSTable{
		Path: path,
		Data: data,
		Keys: make([]string, 0, len(data)),
	}
	for k := range data {
		sst.Keys = append(sst.Keys, k)
	}
	sort.Strings(sst.Keys)
	return sst, nil
}

func (s *SSTable) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.Data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}
