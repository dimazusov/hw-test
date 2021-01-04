package memorystorage

import (
	"context"
	"fmt"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/pkg/errors"
	"sync"
)

type memStorage struct {
	mu sync.Mutex
	events []domain.Event
	lastIndex uint
	storageMaxSize uint
}

type Storage interface {
	Create(ctx context.Context, event domain.Event) (newID uint, err error)
	Update(ctx context.Context, event domain.Event) (err error)
	Delete(ctx context.Context, eventID uint) (err error)
	GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error)
	GetEventsByParams(ctx context.Context, params GettingEventParams) (events []domain.Event, err error)
}

func New(storageMaxSize uint) (Storage, error) {
	return &memStorage{
		mu: sync.Mutex{},
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
	for i, e := range m.events {
		if e.ID == eventID {
			eventCount := len(m.events)

			m.mu.Lock()
			if i == 0 { // front
				if eventCount > 1 {
					m.events = m.events[1:]
				} else {
					m.events = make([]domain.Event, 0)
				}
			} else if i == eventCount - 1 { // last
				m.events = m.events[:i-1]
			} else { // inside
				m.events = append(m.events[:i-1], m.events[i+1:]...)
			}
			m.mu.Unlock()

			return nil
		}
	}

	return errors.Wrap(ErrRecordNotFound, "cannot delete")
}

func (m *memStorage) GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error) {
	for _, e := range m.events {
		if e.ID == eventID {
			return e, nil
		}
	}

	return domain.Event{}, errors.Wrap(ErrRecordNotFound, "cannot found")
}

func (m *memStorage) GetEventsByParams(ctx context.Context, params GettingEventParams) (events []domain.Event, err error) {
	selectedEvents := []domain.Event{}

	for _, e := range m.events {
		if m.isEventHasParams(e, params) {
			selectedEvents = append(selectedEvents, e)
		}
	}

	offset := (params.Page - 1) * params.CountOnPage
	limit := params.CountOnPage

	lastIndex := uint(len(selectedEvents)-1)

	if offset > lastIndex {
		offset = lastIndex
	}

	if limit > lastIndex {
		limit = lastIndex
	}

	return []domain.Event{}, err
}

func (m *memStorage) isEventHasParams(event domain.Event, params GettingEventParams) bool {
	if params.UserID != 0 && event.UserID != params.UserID {
		return false
	}

	if params.ExactTime != 0 {
		if event.Time + uint(event.Timezone) != params.ExactTime + uint(params.Timezone) {
			return false
		}
	}

	if params.FromTime != 0 && event.Time < params.FromTime {
		return false
	}

	if params.ToTime != 0 && event.Time > params.ToTime {
		return false
	}

	if params.UserID != 0 && params.UserID != event.UserID {
		return false
	}

	return true
}