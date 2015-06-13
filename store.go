package adapi

import "errors"

// Store is an interface which when implemented acts as storage backend
type Store interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (interface{}, error)
}

// MemoryStore is a simple in memory implementation of Store. It uses a map.
type MemoryStore struct {
	d map[string]interface{}
}

// NewMemStore initializes a new memoryStore
func NewMemStore() *MemoryStore {
	return &MemoryStore{make(map[string]interface{})}
}

// Set saves the value with the given key
func (m *MemoryStore) Set(key interface{}, value interface{}) error {
	var k string
	switch key.(type) {
	case string:
		k = key.(string)
	default:
		return errors.New("unsupported key type")
	}
	m.d[k] = value
	return nil
}

// Get retrieves a value with the given key.
func (m *MemoryStore) Get(key interface{}) (interface{}, error) {
	var k string
	switch key.(type) {
	case string:
		k = key.(string)
	default:
		return nil, errors.New("unsupported key type")
	}
	if v, ok := m.d[k]; ok {
		return v, nil
	}
	return nil, errors.New("key not found")
}
