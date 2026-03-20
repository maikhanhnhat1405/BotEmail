package store


import (
	"encoding/json"
	"os"
	"sync"
)

type JSONStore struct {
	path string
	mu   sync.RWMutex
	ids  map[string]bool
}

func NewJSONStore(path string) *JSONStore {
	s := &JSONStore{path: path, ids: make(map[string]bool)}
	s.load()
	return s
}

func (s *JSONStore) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ids[id]
}

func (s *JSONStore) Save(id string) error {
	s.mu.Lock()
	s.ids[id] = true
	s.mu.Unlock()
	return s.persist()
}

func (s *JSONStore) load() {
	data, err := os.ReadFile(s.path)
	if err == nil {
		json.Unmarshal(data, &s.ids)
	}
}

func (s *JSONStore) persist() error {
	data, _ := json.Marshal(s.ids)
	return os.WriteFile(s.path, data, 0644)
}