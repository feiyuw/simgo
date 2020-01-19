package storage

import (
	"errors"
	"sync"
)

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{M: map[interface{}]interface{}{}}, nil
}

type MemoryStorage struct {
	sync.RWMutex
	M map[interface{}]interface{}
}

func (ms *MemoryStorage) Add(key interface{}, value interface{}) error {
	ms.M[key] = value
	return nil
}

func (ms *MemoryStorage) Remove(key interface{}) error {
	ms.RLock()
	defer ms.RUnlock()
	if _, exists := ms.M[key]; exists {
		delete(ms.M, key)
	}
	return nil
}

func (ms *MemoryStorage) FindAll() ([]interface{}, error) {
	items := make([]interface{}, len(ms.M))
	idx := 0
	for _, v := range ms.M {
		items[idx] = v
		idx++
	}

	return items, nil
}

func (ms *MemoryStorage) FindOne(key interface{}) (interface{}, error) {
	if v, exists := ms.M[key]; exists {
		return v, nil
	}
	return nil, errors.New("not found")
}
