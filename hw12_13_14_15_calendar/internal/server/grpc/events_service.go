package internalgrpc

import (
	"context"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/pkg/apperror"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/pkg/errors"
)

type EventService struct {
	pb.UnimplementedEventsServer
	app Application
}

func newEventService(app Application) *EventService {
	return &EventService{
		app: app,
	}
}

func (m EventService) GetEventByID(c context.Context, event *pb.Event) (*pb.GetEventRS, error) {
	e, err := m.app.GetEventByID(c, uint(event.Id))
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return nil, errors.Wrap(err, "cannot get event")
	}

	return &pb.GetEventRS{
		Error: err.Error(),
		Event: convertDomainEventToPbEvent(e),
	}, nil
}

func (m EventService) GetEvents(c context.Context, e *pb.Event) (*pb.GetEventsRS, error) {
	events, err := m.app.GetEventsByParams(c, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	pbEvents := []*pb.Event{}
	for _, event := range events {
		pbEvents = append(pbEvents, convertDomainEventToPbEvent(event))
	}

	return &pb.GetEventsRS{
		Error: "",
		Event: pbEvents,
	}, err
}

func (m EventService) UpdateEvent(c context.Context, e *pb.Event) (*pb.UpdateEventRS, error) {
	err := m.app.Update(c, convertPbEventToDomainEvent(e))

	return &pb.UpdateEventRS{}, err
}

func (m EventService) CreateEvent(c context.Context, e *pb.Event) (*pb.CreateEventRS, error) {
	newID, err := m.app.Create(c, convertPbEventToDomainEvent(e))

	return &pb.CreateEventRS{
		Event: uint32(newID),
	}, err
}

func (m EventService) DeleteEvent(c context.Context, e *pb.Event) (*pb.DeleteEventRS, error) {
	err := m.app.Delete(c, uint(e.Id))

	return &pb.DeleteEventRS{}, err
}

func convertDomainEventToPbEvent(event domain.Event) *pb.Event {
	return &pb.Event{
		Id:               uint32(event.ID),
		Title:            event.Title,
		Time:             uint32(event.Time),
		Timezone:         uint32(event.Timezone),
		Duration:         uint32(event.Duration),
		Description:      event.Description,
		UserId:           uint32(event.UserID),
		NotificationTime: uint32(event.NotificationTime),
	}
}

func convertPbEventToDomainEvent(event *pb.Event) domain.Event {
	return domain.Event{
		ID:               uint(event.Id),
		Title:            event.Title,
		Time:             uint(event.Time),
		Timezone:         uint8(event.Timezone),
		Duration:         uint(event.Duration),
		Description:      event.Description,
		UserID:           uint(event.UserId),
		NotificationTime: uint(event.NotificationTime),
	}
}
