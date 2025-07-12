package cache

import (
	"sync"
	"time"
)

func NewInmem[V any]() *Inmem[V] {
	return &Inmem[V]{
		m:    &sync.RWMutex{},
		data: map[string]*inmemItem[V]{},
	}
}

type inmemItem[T any] struct {
	Value      T
	Expiration time.Time
}

type Inmem[V any] struct {
	m    *sync.RWMutex
	data map[string]*inmemItem[V]
}

func (c *Inmem[V]) Get(key string) (V, bool) {
	c.m.RLock()
	defer c.m.RUnlock()

	var v V

	item, ok := c.data[key]
	if !ok {
		return v, false
	}
	if time.Now().Before(item.Expiration) {
		return item.Value, true
	}

	return v, false
}

func (c *Inmem[V]) GetEx(key string, ttl time.Duration) (V, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	item, ok := c.data[key]
	if !ok {
		var v V
		return v, false
	}

	if time.Now().Before(item.Expiration) {
		c.data[key].Expiration = time.Now().Add(ttl)
		return item.Value, true
	}

	delete(c.data, key)

	var v V
	return v, false
}

func (c *Inmem[V]) Set(key string, value V) {
	c.m.Lock()
	defer c.m.Unlock()

	item := &inmemItem[V]{
		Value:      value,
		Expiration: time.Now().Add(24 * time.Hour),
	}
	c.data[key] = item
}

func (c *Inmem[V]) SetEx(key string, value V, ttl time.Duration) {
	c.m.Lock()
	defer c.m.Unlock()

	item := &inmemItem[V]{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
	c.data[key] = item
}

func (c *Inmem[V]) Del(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	delete(c.data, key)
}
