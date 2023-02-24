package cache

import (
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, value []byte, timeToLive time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
}

func (c *Cache) Has(key []byte) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(key))

	return nil
}
