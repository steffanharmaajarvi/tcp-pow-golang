package storage

import (
	"sync"
	"time"
)

type InMemoryStorage struct {
	dataMap map[string]inMemoryValue
	lock    *sync.Mutex
	time    time.Time
}

type inMemoryValue struct {
	SetTime    int64
	Expiration int64
}

func InitInMemoryStorage(time time.Time) *InMemoryStorage {
	return &InMemoryStorage{
		dataMap: make(map[string]inMemoryValue, 0),
		lock:    &sync.Mutex{},
		time:    time,
	}
}

// Add - add value with expiration (in seconds) to storage
func (c *InMemoryStorage) Add(key string, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataMap[key] = inMemoryValue{
		SetTime:    c.time.Unix(),
		Expiration: expiration,
	}
	return nil
}

// Get - check existence of int key in storage
func (c *InMemoryStorage) Get(key string) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.dataMap[key]
	if ok && c.time.Unix()-value.SetTime > value.Expiration {
		return false, nil
	}
	return ok, nil
}

// Delete - delete key from storage
func (c *InMemoryStorage) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.dataMap, key)
}
