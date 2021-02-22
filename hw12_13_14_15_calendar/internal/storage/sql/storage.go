//nolint:golint
package sql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/jmoiron/sqlx"
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
	DeleteOldEvents(ctx context.Context, timeTo uint) (err error)
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
	timestamp := toTimestamp(event.Time)
	notifTimestamp := toTimestamp(event.NotificationTime)

	params := []interface{}{
		event.Title,
		timestamp,
		event.Timezone,
		event.Duration,
		event.Description,
		event.UserID,
		notifTimestamp,
	}

	query := "INSERT INTO event (title, time, timezone, duration, description, user_id, notification_time) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	err = m.conn.QueryRowContext(ctx, query, params...).Scan(&newID)
	if err != nil {
		return 0, errors.Wrap(err, "cannot create event")
	}

	return newID, nil
}

func (m *postgresStorage) Update(ctx context.Context, event domain.Event) (err error) {
	timestamp := toTimestamp(event.Time)
	notifTimestamp := toTimestamp(event.NotificationTime)

	params := []interface{}{
		event.Title,
		timestamp,
		event.Timezone,
		event.Duration,
		event.Description,
		event.UserID,
		notifTimestamp,
		event.IsNotificationSend,
		event.ID,
	}
	query := "UPDATE event SET title=$1," +
		"time=$2," +
		"timezone=$3," +
		"duration=$4," +
		"description=$5," +
		"user_id=$6," +
		"notification_time=$7," +
		"is_notification_send=$8 " +
		"where id = $9"

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

func (m *postgresStorage) GetEventByID(ctx context.Context, eventID uint) (e domain.Event, err error) {
	query := "select id, title, time, timezone, duration, description, user_id, notification_time, is_notification_send " +
		"FROM event WHERE id = $1"

	row := m.conn.QueryRowxContext(ctx, query, eventID)

	var timestamp string
	var notifTimestamp string

	err = row.Scan(&e.ID, &e.Title, &timestamp, &e.Timezone, &e.Duration, &e.Description, &e.UserID, &notifTimestamp, &e.IsNotificationSend)

	e.Time = toUnixTime(timestamp)
	e.NotificationTime = toUnixTime(notifTimestamp)

	if errors.Is(err, sql.ErrNoRows) {
		return e, errors.Wrap(ErrRecordNotFound, "cannot get event")
	}

	return e, nil
}

func (m *postgresStorage) GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error) {
	query := "select id, title, time, timezone, duration, description, user_id, notification_time, is_notification_send FROM event "
	where := []string{}
	qParams := make(map[string]interface{})

	userID, ok := params["userID"].(uint)
	if ok {
		where = append(where, " user_id = :user_id")
		qParams["user_id"] = userID
	}

	title, ok := params["userID"].(string)
	if ok {
		where = append(where, " title = :title")
		qParams["title"] = title
	}

	ids, ok := params["ids"]
	if ok {
		where = append(where, " id IN (:ids)")
		qParams["ids"] = ids
	}

	notifTimeFrom, ok := params["notificationTimeFrom"].(uint)
	if ok {
		tm := time.Unix(int64(notifTimeFrom), 0)
		formatTime := tm.Format(time.RFC3339)

		where = append(where, " notification_time < :notification_time_from")
		qParams["notification_time_from"] = formatTime
	}

	isNotificationSend, ok := params["isNotificationSend"].(string)
	if ok {
		where = append(where, " is_notification_send = :is_notification_send")
		qParams["is_notification_send"] = isNotificationSend
	}

	if len(where) != 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	countOnPage, ok := params["countOnPage"].(uint)
	if ok {
		query += " LIMIT :countOnPage"
		qParams["count_on_page"] = countOnPage
	}

	page, ok := params["page"].(uint)
	if ok {
		query += " OFFSET :offset"
		qParams["offset"] = (page - 1) * countOnPage
	}

	rows, err := m.conn.NamedQueryContext(ctx, query, qParams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp string
		var timestampNotification string

		var e domain.Event
		err := rows.Scan(&e.ID, &e.Title, &timestamp, &e.Timezone, &e.Duration, &e.Description, &e.UserID, &timestampNotification, &e.IsNotificationSend)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get events")
		}

		e.Time = toUnixTime(timestamp)
		e.NotificationTime = toUnixTime(timestampNotification)

		events = append(events, e)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "cannot get events")
	}

	return events, nil
}

func (m *postgresStorage) DeleteOldEvents(ctx context.Context, timeTo uint) (err error) {
	d, err := time.ParseDuration(strconv.Itoa(int(timeTo)) + "s")
	if err != nil {
		return err
	}

	expiredTime := time.Now().Add(-d).Unix()

	expiredTimestamp := toTimestamp(uint(expiredTime))

	_, err = m.conn.ExecContext(ctx, "DELETE FROM event WHERE time < $1", expiredTimestamp)
	if err != nil {
		return err
	}

	return nil
}

func toUnixTime(t string) uint {
	parsedTime, _ := time.Parse(time.RFC3339, t)

	return uint(parsedTime.Unix())
}

func toTimestamp(unixTime uint) string {
	return time.Unix(int64(unixTime), 0).Format(time.RFC3339)
}
