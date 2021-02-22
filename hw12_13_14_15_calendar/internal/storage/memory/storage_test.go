//nolint:ineffassign
package memorystorage

import (
	"context"
	"errors"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStorage_GetEventByID(t *testing.T) {
	repository, err := New(20)
	require.Nil(t, err)

	event := domain.Event{
		Title:            "title",
		Time:             10,
		Timezone:         0,
		Duration:         0,
		Description:      "desc",
		UserID:           23,
		NotificationTime: 9,
	}

	newID, err := repository.Create(context.Background(), event)
	require.Nil(t, err)
	require.Equal(t, uint(1), newID)

	event, err = repository.GetEventByID(context.Background(), newID)
	require.Nil(t, err)
	require.Equal(t, "title", event.Title)
}

func TestMemStorage_Create(t *testing.T) {
	repository, err := New(20)
	require.Nil(t, err)

	event := domain.Event{
		Title:            "title",
		Time:             10,
		Timezone:         0,
		Duration:         0,
		Description:      "desc",
		UserID:           23,
		NotificationTime: 9,
	}

	newID, err := repository.Create(context.Background(), event)
	require.Nil(t, err)
	require.Equal(t, uint(1), newID)
}

func TestMemStorage_Delete(t *testing.T) {
	repository, err := New(20)
	require.Nil(t, err)

	event := domain.Event{
		Title:            "title",
		Time:             10,
		Timezone:         0,
		Duration:         0,
		Description:      "desc",
		UserID:           23,
		NotificationTime: 9,
	}

	newID, err := repository.Create(context.Background(), event)
	require.Nil(t, err)
	require.Equal(t, uint(1), newID)

	err = repository.Delete(context.Background(), newID)
	require.Nil(t, err)

	event, err = repository.GetEventByID(context.Background(), newID)
	require.True(t, errors.Is(err, ErrRecordNotFound))
}

func TestMemStorage_Update(t *testing.T) {
	repository, err := New(20)
	require.Nil(t, err)

	event := domain.Event{
		Title:            "title",
		Time:             10,
		Timezone:         0,
		Duration:         0,
		Description:      "desc",
		UserID:           23,
		NotificationTime: 9,
	}

	newID, err := repository.Create(context.Background(), event)
	require.Nil(t, err)
	require.Equal(t, uint(1), newID)

	event.ID = newID
	event.Title = "event title"
	event.Description = "event description"

	err = repository.Update(context.Background(), event)
	require.Nil(t, err)

	event, err = repository.GetEventByID(context.Background(), newID)
	require.Nil(t, err)

	require.Equal(t, "event title", event.Title)
	require.Equal(t, "event description", event.Description)
}
