package memorystorage

import (
	"context"
	"sync"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/pkg/errors"
)

type memStorage struct {
	mu             sync.Mutex
	events         []domain.Event
	lastIndex      uint
	storageMaxSize uint
}

type Storage interface {
	Create(ctx context.Context, event domain.Event) (newID uint, err error)
	Update(ctx context.Context, event domain.Event) (err error)
	Delete(ctx context.Context, eventID uint) (err error)
	GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error)
	GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error)
}

func New(storageMaxSize uint) (Storage, error) {
	return &memStorage{
		mu:             sync.Mutex{},
		storageMaxSize: storageMaxSize,
	}, nil
}

func (m *memStorage) Create(ctx context.Context, event domain.Event) (newID uint, err error) {
	m.mu.Lock()
	m.lastIndex++
	event.ID = m.lastIndex
	m.events = append(m.events, event)
	m.mu.Unlock()

	return event.ID, nil
}

func (m *memStorage) Update(ctx context.Context, event domain.Event) (err error) {
	for i, e := range m.events {
		if e.ID == event.ID {
			m.mu.Lock()
			m.events[i] = event
			m.mu.Unlock()
			break
		}
	}

	return nil
}

func (m *memStorage) Delete(ctx context.Context, eventID uint) (err error) {
	index, ok := m.findIndexByEventID(eventID)
	if !ok {
		return errors.Wrap(ErrRecordNotFound, "cannot delete")
	}

	eventCount := len(m.events)
	m.mu.Lock()
	switch {
	case index == 0: // front
		if eventCount > 1 {
			m.events = m.events[1:]
		} else {
			m.events = make([]domain.Event, 0)
		}
	case index == eventCount-1: // last
		m.events = m.events[:index-1]
	default: // inside
		m.events = append(m.events[:index-1], m.events[index+1:]...)
	}
	m.mu.Unlock()

	return nil
}

func (m *memStorage) findIndexByEventID(eventID uint) (int, bool) {
	for i, e := range m.events {
		if e.ID == eventID {
			return i, true
		}
	}

	return 0, false
}

func (m *memStorage) GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error) {
	for _, e := range m.events {
		if e.ID == eventID {
			return e, nil
		}
	}

	return domain.Event{}, errors.Wrap(ErrRecordNotFound, "cannot found")
}

func (m *memStorage) GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error) {
	selectedEvents := []domain.Event{}

	for _, e := range m.events {
		if m.isEventHasParams(e, params) {
			selectedEvents = append(selectedEvents, e)
		}
	}

	page := 1
	countOnPage := 10

	if _, ok := params["page"]; ok {
		page = params["page"].(int)
	}

	if _, ok := params["countOnPage"]; ok {
		countOnPage = params["countOnPage"].(int)
	}

	offset := (page - 1) * countOnPage
	limit := countOnPage

	lastIndex := len(selectedEvents) - 1

	if offset > lastIndex {
		offset = lastIndex
	}

	if limit > lastIndex {
		limit = lastIndex
	}

	return selectedEvents[offset:limit], err
}

func (m *memStorage) isEventHasParams(event domain.Event, params map[string]interface{}) bool {
	_, ok := params["userID"]
	if ok && event.UserID != params["userID"].(uint) {
		return false
	}

	_, ok = params["exactTime"]
	if ok && params["exactTime"] != 0 {
		var timezone uint
		timezone, _ = params["timezone"].(uint)

		if event.Time+uint(event.Timezone) != params["exactTime"].(uint)+timezone {
			return false
		}
	}

	fromTime, ok := params["fromTime"].(uint)
	if ok && event.Time < fromTime {
		return false
	}

	toTime, ok := params["toTime"].(uint)
	if ok && event.Time > toTime {
		return false
	}

	return true
}
