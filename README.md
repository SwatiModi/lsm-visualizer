# LSM Tree Visualizer

This project simulates a basic Log-Structured Merge Tree (LSM Tree) data structure in Go, focusing on core components like Memtable, Write-Ahead Log (WAL), SSTables, and compaction. It’s designed as a learning and visualization tool to demonstrate how modern databases like RocksDB, LevelDB, and Cassandra manage data efficiently.

---

## Features

- **Memtable:** In-memory sorted key-value store
- **WAL:** Write-ahead log for durability
- **SSTable:** Immutable on-disk sorted string tables saved as JSON files
- **Compaction:** Merge multiple SSTables into a single one, simulating database compaction
- **Bloom Filter:** Simple probabilistic data structure to speed up reads by filtering non-existent keys
- **Interactive CLI:** Insert (`put`), retrieve (`get`), compact SSTables, and exit commands

---

## Why LSM Trees?

LSM Trees optimize write throughput by buffering writes in memory (memtable) and flushing them to disk in sorted files (SSTables). They perform periodic compactions to merge and keep SSTables optimized for reads. They are widely used in many high-performance databases.

---

## Code Structure

- `lsm/memtable.go`: Memtable implementation — a sorted in-memory key-value map with concurrency safety.
- `lsm/wal.go`: Write-Ahead Log — an append-only file logging all writes for durability.
- `lsm/sstable.go`: SSTable simulation — saves immutable sorted key-value maps to disk as JSON files.
- `lsm/compactor.go`: Compaction logic — merges two SSTables into one.
- `lsm/bloom.go`: Simple Bloom filter for membership checks.
- `lsm/store.go`: The main LSM Store managing memtable, WAL, SSTables, and compaction.
- `main.go`: Command-line interface to interact with the LSM Store.

---

## How to Run

1. Clone the repo and initialize Go modules:

   ```bash
   go mod tidy
