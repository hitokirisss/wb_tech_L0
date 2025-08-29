package cache

import (
	"encoding/json"
	"sync"
)

type Cache struct {
	mu sync.RWMutex
	m  map[string]json.RawMessage
}

func New() *Cache { return &Cache{m: make(map[string]json.RawMessage)} }

func (c *Cache) Get(id string) (json.RawMessage, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[id]
	return v, ok
}

func (c *Cache) Set(id string, raw json.RawMessage) {
	c.mu.Lock()
	c.m[id] = raw
	c.mu.Unlock()
}
