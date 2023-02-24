package cache

import "sync"

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}
