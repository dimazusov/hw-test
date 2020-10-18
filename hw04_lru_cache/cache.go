package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]listItem
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]listItem),
	}
}

func (m *lruCache) Set(key Key, value interface{}) bool {
	item, wasInCache := m.items[key]
	if wasInCache {
		cItem := item.Value.(cacheItem)
		cItem.Value = value
		item.Value = cItem
		m.items[key] = item

		return true
	}

	if m.queue.Len() == m.capacity {
		last := m.queue.Back()
		delete(m.items, last.Value.(cacheItem).Key)
		m.queue.Remove(last)
	}

	lItem := m.queue.PushFront(cacheItem{
		Key:   key,
		Value: value,
	})
	m.items[key] = *lItem

	return false
}

func (m *lruCache) Get(key Key) (interface{}, bool) {
	cItem, ok := m.items[key]
	if !ok {
		return nil, false
	}

	return cItem.Value.(cacheItem).Value, true
}

func (m *lruCache) Clear() {
	for i := range m.items {
		delete(m.items, i)
	}
}
