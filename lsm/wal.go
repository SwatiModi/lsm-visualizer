package lsm

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type WAL struct {
	mu   sync.Mutex
	file *os.File
	w    *bufio.Writer
}

func NewWAL(filename string) (*WAL, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{
		file: f,
		w:    bufio.NewWriter(f),
	}, nil
}

func (w *WAL) Write(key, value string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := fmt.Fprintf(w.w, "%s=%s\n", key, value)
	if err != nil {
		return err
	}
	return w.w.Flush()
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}
