package sql

import (
	"context"
	"database/sql"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/jmoiron/sqlx"
	// comment
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type postgresStorage struct {
	conn       *sqlx.DB
	dsn        string
	driverName string
}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	Create(ctx context.Context, event domain.Event) (newID uint, err error)
	Update(ctx context.Context, event domain.Event) (err error)
	Delete(ctx context.Context, eventID uint) (err error)
	GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error)
	GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error)
}

func New(driverName, dsn string) (Storage, error) {
	conn, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}

	return &postgresStorage{
		conn:       conn,
		dsn:        dsn,
		driverName: driverName,
	}, nil
}

func (m *postgresStorage) Connect(ctx context.Context) error {
	err := m.conn.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}

	return nil
}

func (m *postgresStorage) Close() error {
	err := m.conn.Close()
	if err != nil {
		return errors.Wrap(err, "cannot close")
	}

	return nil
}

func (m *postgresStorage) Create(ctx context.Context, event domain.Event) (newID uint, err error) {
	params := []interface{}{
		event.Title,
		event.Time,
		event.Timezone,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationTime,
	}

	query := "INSERT INTO event (title, time, timezone, duration, description, user_id, notification_time) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7);"

	res, err := m.conn.ExecContext(ctx, query, params...)
	if err != nil {
		return 0, errors.Wrap(err, "cannot create event")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "cannot get last id")
	}

	return uint(lastID), nil
}

func (m *postgresStorage) Update(ctx context.Context, event domain.Event) (err error) {
	params := []interface{}{
		event.Title,
		event.Time,
		event.Timezone,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationTime,
		event.ID,
	}
	query := "UPDATE event SET title=$1," +
		"time=$2," +
		"timezone=$3," +
		"duration=$4," +
		"description=$5," +
		"user_id=$6," +
		"notification_time=$7 " +
		"where id = $8"

	res, err := m.conn.ExecContext(ctx, query, params...)
	if err != nil {
		return errors.Wrap(err, "cannot update event")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get update rows affected")
	}

	if affected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *postgresStorage) Delete(ctx context.Context, eventID uint) (err error) {
	query := "DELETE FROM event where id = $1"

	res, err := m.conn.ExecContext(ctx, query, eventID)
	if err != nil {
		return errors.Wrap(err, "cannot delete event")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "cannot get delete rows affected")
	}

	if affected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *postgresStorage) GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error) {
	query := "select id, title, time, timezone, duration, description, user_id, notification_time " +
		"FROM event WHERE id = $1"

	err = m.conn.GetContext(ctx, &event, query, eventID)
	if errors.Is(err, sql.ErrNoRows) {
		return event, errors.Wrap(ErrRecordNotFound, "cannot get event")
	}

	return event, nil
}

func (m *postgresStorage) GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error) {
	query := "select id, title, time, timezone, duration, description, user_id, notification_time " +
		"FROM event WHERE 1 "
	qParams := []interface{}{}

	userID, ok := params["userID"].(uint)
	if ok {
		query += " user_id = ? "
		qParams = append(qParams, userID)
	}

	title, ok := params["userID"].(string)
	if ok {
		query += " title = ? "
		qParams = append(qParams, title)
	}

	ids, ok := params["ids"]
	if ok {
		query += " id IN (?) "
		qParams = append(qParams, ids)
	}

	countOnPage, ok := params["countOnPage"].(uint)
	if ok {
		query += " LIMIT ?"
		qParams = append(qParams, countOnPage)
	}

	page, ok := params["page"].(uint)
	if ok {
		query += " OFFSET ?"
		qParams = append(qParams, (page-1)*countOnPage)
	}

	rows, _ := m.conn.QueryContext(ctx, query, qParams...)
	defer rows.Close()

	for rows.Next() {
		var e domain.Event
		err := rows.Scan(&e.ID, &e.Title, &e.Time, &e.Timezone, &e.Duration, &e.Description, &e.UserID, &e.NotificationTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get events")
		}
		events = append(events, e)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "cannot get events")
	}

	return events, nil
}
