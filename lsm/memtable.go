package lsm

type Memtable struct {
	data     map[string]string
	capacity int
}

func NewMemtable(capacity int) *Memtable {
	return &Memtable{
		data:     make(map[string]string),
		capacity: capacity,
	}
}

func (m *Memtable) Put(key, val string) {
	m.data[key] = val
}

func (m *Memtable) Get(key string) (string, bool) {
	val, found := m.data[key]
	return val, found
}

func (m *Memtable) Size() int {
	return len(m.data)
}

func (m *Memtable) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *Memtable) Flush() map[string]string {
	data := m.data
	m.data = make(map[string]string)
	return data
}
