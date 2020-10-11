package hw04_lru_cache //nolint:golint,stylecheck

func (m list) Len() int {
	return m.len
}

func (m *list) Front() *listItem {
	if m.len == 0 {
		return nil
	}

	return m.First
}

func (m list) Back() *listItem {
	if m.len == 0 {
		return nil
	}

	return m.Last
}

func (m *list) PushFront(v interface{}) *listItem {
	m.len++
	lItem := &listItem{
		Value: v,
		Prev:  nil,
	}

	if m.len == 1 {
		m.First = lItem
		m.Last = lItem
		return lItem
	}

	lItem.Next = m.First
	m.First.Prev = lItem
	m.First = lItem

	return lItem
}

func (m *list) PushBack(v interface{}) *listItem {
	m.len++
	lItem := &listItem{
		Value: v,
		Next:  nil,
	}

	if m.len == 1 {
		m.First = lItem
		m.Last = lItem
		return lItem
	}

	lItem.Prev = m.Last
	m.Last.Next = lItem
	m.Last = lItem

	return lItem
}

func (m *list) Remove(lItem *listItem) {
	m.len--

	if lItem.Prev == nil && lItem.Next == nil {
		m.First = nil
		m.Last = nil
		return
	}

	if lItem.Prev == nil && lItem.Next != nil {
		m.First = lItem.Next
		lItem.Next.Prev = nil
		return
	}

	if lItem.Next == nil && lItem.Prev != nil {
		m.Last = lItem.Prev
		lItem.Prev.Next = nil
		return
	}

	prevItem := lItem.Prev
	nextItem := lItem.Next

	prevItem.Next = nextItem
	nextItem.Prev = prevItem
}

func (m *list) MoveToFront(l *listItem) {
	m.Remove(l)
	m.PushFront(l.Value)
}
