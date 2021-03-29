package hw04lrucache

import (
	"sync"
)

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*Element
	lock     sync.RWMutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	e, ok := l.items[key]
	if ok {
		e.Value = cacheItem{key, value}
		l.items[key] = e
		l.queue.MoveToFront(e)
	}
	e = l.queue.PushFront(cacheItem{key, value})
	l.items[key] = e

	if l.queue.Len() > l.capacity {
		l.queue.Remove(l.queue.Back())
		delete(l.items, l.queue.Back().Value.(cacheItem).key)
	}
	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if e, ok := l.items[key]; ok {
		return e.Value.(cacheItem).value, ok
	} else {
		return nil, ok
	}
}

func (l *lruCache) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*Element, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Element, capacity),
	}
}
