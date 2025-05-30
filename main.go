package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/swatimodi/lsmvisualizer/lsm"
)

var (
	tree    *lsm.LSMTree
	treeMux sync.RWMutex
)

func main() {
	tree = lsm.NewLSMTree(1000, 20) // memtable capacity, bloom filter size
	os.MkdirAll("sstables", os.ModePerm)

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/memtable", handleMemtable)
	http.HandleFunc("/bloom", handleBloom)
	http.HandleFunc("/sstables", handleSSTables)
	http.HandleFunc("/compactions", handleCompactions)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ui/index.html")
}

type putRequest struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req putRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Key == "" || req.Val == "" {
		http.Error(w, "missing key or val", http.StatusBadRequest)
		return
	}

	treeMux.Lock()
	tree.Put(req.Key, req.Val)
	treeMux.Unlock()

	fmt.Fprintf(w, "inserted %s -> %s\n", req.Key, req.Val)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	treeMux.RLock()
	val, found := tree.Get(key)
	treeMux.RUnlock()

	if !found {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", val)
}

func handleMemtable(w http.ResponseWriter, r *http.Request) {
	treeMux.RLock()
	defer treeMux.RUnlock()

	keys := tree.MemtableKeys()
	json.NewEncoder(w).Encode(keys)
}

func handleBloom(w http.ResponseWriter, r *http.Request) {
	treeMux.RLock()
	defer treeMux.RUnlock()

	b := tree.BloomStats()
	json.NewEncoder(w).Encode(b)
}

func handleSSTables(w http.ResponseWriter, r *http.Request) {
	treeMux.RLock()
	defer treeMux.RUnlock()

	levels := tree.SSTablesMetadata()
	json.NewEncoder(w).Encode(levels)
}

func handleCompactions(w http.ResponseWriter, r *http.Request) {
	treeMux.RLock()
	defer treeMux.RUnlock()

	logs := tree.CompactionLogs()
	json.NewEncoder(w).Encode(logs)
}
