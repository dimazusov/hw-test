package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(l *listItem)
	MoveToFront(l *listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	len   int
	First *listItem
	Last  *listItem
}

func NewList() List {
	return &list{}
}
