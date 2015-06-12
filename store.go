package adapi

import "errors"

type Store interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (interface{}, error)
}

type MemoryStore struct {
	d map[string]interface{}
}

func NewMemSTore() *MemoryStore {
	return &MemoryStore{make(map[string]interface{})}
}

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
