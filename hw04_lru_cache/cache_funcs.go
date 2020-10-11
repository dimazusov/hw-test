package hw04_lru_cache

func (m *lruCache) Set(key string, value interface{}) bool {
	cItem := cacheItem{
		Key: key,
		Value: value,
	}

	var lItem *listItem
	_, wasInCache := m.items[key]
	if wasInCache {
		for el := m.queue.Front() ; el != nil; el = el.Next {
			if el.Value.(cacheItem).Key == key {
				lItem = el
				lItem.Value = cacheItem{
					Key: key,
					Value: value,
				}

				m.queue.MoveToFront(el)
				break
			}
		}
	} else {
		if m.queue.Len() == m.capacity {
			last := m.queue.Back()
			delete(m.items, last.Value.(cacheItem).Key)
			m.queue.Remove(last)
		}
		lItem = m.queue.PushFront(cItem)
	}

	m.items[key] = *lItem

	return wasInCache
}

func (m *lruCache) Get(key string) (interface{}, bool) {
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