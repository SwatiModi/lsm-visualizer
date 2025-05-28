package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type LSMTree struct {
	memtable map[string]string
	sstables []map[string]string
	mu       sync.RWMutex
}

func NewLSMTree() *LSMTree {
	return &LSMTree{
		memtable: make(map[string]string),
		sstables: []map[string]string{},
	}
}

func (db *LSMTree) Put(key, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.memtable[key] = value

	if len(db.memtable) >= 5 {
		db.flushMemtable()
	}
	return nil
}

func (db *LSMTree) Get(key string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if val, ok := db.memtable[key]; ok {
		return val, nil
	}

	for i := len(db.sstables) - 1; i >= 0; i-- {
		if val, ok := db.sstables[i][key]; ok {
			return val, nil
		}
	}
	return "", errors.New("key not found")
}

func (db *LSMTree) flushMemtable() {
	sstableCopy := make(map[string]string)
	for k, v := range db.memtable {
		sstableCopy[k] = v
	}
	db.sstables = append(db.sstables, sstableCopy)
	db.memtable = make(map[string]string)
}

func (db *LSMTree) Compact() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	n := len(db.sstables)
	if n < 2 {
		return errors.New("not enough SSTables to compact")
	}

	sstable1 := db.sstables[n-2]
	sstable2 := db.sstables[n-1]
	merged := make(map[string]string)

	for k, v := range sstable1 {
		merged[k] = v
	}
	for k, v := range sstable2 {
		merged[k] = v
	}

	db.sstables = append(db.sstables[:n-2], merged)
	return nil
}

func (db *LSMTree) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data := map[string]interface{}{
		"memtable": db.memtable,
		"sstables": db.sstables,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	db := NewLSMTree()

	// Start HTTP server for frontend visualization
	go func() {
		fmt.Println("Starting HTTP server at http://localhost:8080/status")
		http.Handle("/status", db)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("HTTP server error:", err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("LSM Tree CLI - commands: put key value | get key | compact | exit")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]

		switch cmd {
		case "put":
			if len(parts) < 3 {
				fmt.Println("Usage: put key value")
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			if err := db.Put(key, value); err != nil {
				fmt.Println("Put error:", err)
			} else {
				fmt.Println("Inserted", key)
			}

		case "get":
			if len(parts) != 2 {
				fmt.Println("Usage: get key")
				continue
			}
			key := parts[1]
			val, err := db.Get(key)
			if err != nil {
				fmt.Println("Get error:", err)
			} else {
				fmt.Println(key, "=", val)
			}

		case "compact":
			err := db.Compact()
			if err != nil {
				fmt.Println("Compact error:", err)
			} else {
				fmt.Println("Compaction done")
			}

		case "exit":
			fmt.Println("Bye!")
			return

		default:
			fmt.Println("Unknown command:", cmd)
			fmt.Println("Commands: put key value | get key | compact | exit")
		}
	}
}
