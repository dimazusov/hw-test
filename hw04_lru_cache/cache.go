package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue List
	items map[string]listItem
}

type cacheItem struct {
	Key string
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue: NewList(),
		items: make(map[string]listItem),
	}
}
