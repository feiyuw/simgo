package storage

import (
	"errors"
	"sync"
)

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{M: map[string]interface{}{}}, nil
}

type MemoryStorage struct {
	sync.RWMutex
	M map[string]interface{}
}

func (ms *MemoryStorage) Add(key string, value interface{}) error {
	ms.Lock()
	defer ms.Unlock()
	ms.M[key] = value
	return nil
}

func (ms *MemoryStorage) Remove(key string) error {
	ms.Lock()
	defer ms.Unlock()
	if _, exists := ms.M[key]; exists {
		delete(ms.M, key)
	}
	return nil
}

func (ms *MemoryStorage) FindAll() ([]interface{}, error) {
	ms.RLock()
	defer ms.RUnlock()
	items := make([]interface{}, len(ms.M))
	idx := 0
	for _, v := range ms.M {
		items[idx] = v
		idx++
	}

	return items, nil
}

func (ms *MemoryStorage) FindOne(key string) (interface{}, error) {
	ms.RLock()
	defer ms.RUnlock()
	if v, exists := ms.M[key]; exists {
		return v, nil
	}
	return nil, errors.New("not found")
}
