package cachex

import (
	"sync"
	"time"
)

type ExpireFunc func()

type cacheItem struct {
	value   interface{}
	addTime time.Time
	auto    func() interface{}
}

type Cache struct {
	m          map[interface{}]*cacheItem
	lock       sync.RWMutex
	expiration time.Duration
}

func NewCache(expiration time.Duration) *Cache {
	cache := Cache{
		m:          make(map[interface{}]*cacheItem),
		expiration: expiration,
		lock:       sync.RWMutex{},
	}
	return &cache
}

func (c *Cache) Set(key interface{}, value interface{}) ExpireFunc {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.m[key] = &cacheItem{
		value:   value,
		addTime: time.Now(),
		auto:    nil,
	}
	return func() {
		c.Expire(key)
	}
}

func (c *Cache) SetAndAutoFunc(key interface{}, value interface{}, auto func() interface{}) ExpireFunc {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.m[key] = &cacheItem{
		value:   value,
		addTime: time.Now(),
		auto:    auto,
	}
	return func() {
		c.Expire(key)
	}
}

func (c *Cache) Get(key interface{}) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if item, ok := c.m[key]; ok {
		if item == nil {
			return nil
		}
		if c.expiration > time.Now().Sub(item.addTime) {
			return item.value
		} else {
			c.Set(key, nil)
			if item.auto != nil {
				v := item.auto()
				c.Set(key, v)
				return v
			} else {
				return nil
			}
		}
	} else {
		return nil
	}

}

func (c *Cache) Expire(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.m[key] = nil
	delete(c.m, key)
}
