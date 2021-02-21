package integration_tests

import (
	"github.com/jmoiron/sqlx"
	"os"
	"log"
	"testing"

	"github.com/pressly/goose"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
)

const testConfigPath = "../../configs/config_test.yaml"
const migrationDir = "../../migrations"

func RunTests(t *testing.T) {
	cfg, err := config.New(testConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg.DB.Postgres.Dialect)
	log.Println(cfg.DB.Postgres.Dsn)

	conn, err := sqlx.Connect(cfg.DB.Postgres.Dialect, cfg.DB.Postgres.Dsn)
	require.Nil(t, err)


	err = goose.Up(conn.DB, migrationDir)
	require.Nil(t, err)

	os.Exit(0)






	//conn, err := getConnDb(dsn)
	//require.Nil(t, err)
	//

	//
	//cfg, err := config.New(testConfigPath)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//cfg.DB.Postgres.Dsn = dsn
	//
	//lg, err := logger.New(cfg.Logger.Path, cfg.Logger.Level)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//rep, err := storage.NewRepository(cfg)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//calendar := app.New(lg, rep.(app.Repository))
	//require.NotNil(t, calendar)
	//
	//params := map[string]interface{}{}
	//
	//events, err := calendar.GetEventsByParams(context.Background(), params)
	//require.Nil(t, err)
	//require.Equal(t, 0, len(events))
	//
	//srv := internalhttp.NewServer(cfg, calendar)
	//require.NotNil(t, srv)
	//
	//go func() {
	//	srv.Start(context.Background())
	//}()
	//
	//time.Sleep(2 * time.Second)
	//
	//event := domain.Event{
	//	Title:              "title",
	//	Time:               1713503896,
	//	Timezone:           3,
	//	Duration:           0,
	//	Description:        "test",
	//	UserID:             5,
	//	NotificationTime:   1713503796,
	//	IsNotificationSend: false,
	//}
	//
	//event.ID = testCreateEvent(t, cfg, event)
	//testGetEvents(t, cfg, event)
	//testGetEventByID(t, cfg, event)
	//
	//srv.Stop(context.Background())
}

//func testCreateEvent(t *testing.T, cfg *config.Config, event domain.Event) uint {
//	b, err := json.Marshal(event)
//	require.Nil(t, err)
//	require.NotNil(t, b)
//
//	url := fmt.Sprintf("http://%s:%s/event", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
//	response, statusCode := doRequest(t, http.MethodPost, url, b)
//	require.Equal(t, http.StatusOK, statusCode)
//
//	res := struct {
//		NewID uint `json:"newID"`
//	}{}
//	err = json.Unmarshal(response, &res)
//	require.Nil(t, err)
//	require.NotEqual(t, 0, res.NewID)
//
//	return res.NewID
//}
//
//func testGetEvents(t *testing.T, cfg *config.Config, createdEvent domain.Event) {
//	url := fmt.Sprintf("http://%s:%s/events", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
//	b, statusCode := doRequest(t, http.MethodGet, url, []byte("{}"))
//	require.Equal(t, http.StatusOK, statusCode)
//
//	events := []domain.Event{}
//	err := json.Unmarshal(b, &events)
//	require.Nil(t, err)
//	require.Equal(t, len(events), 1)
//
//	e := events[0]
//	require.Equal(t, e.UserID, createdEvent.UserID)
//	require.Equal(t, e.IsNotificationSend, createdEvent.IsNotificationSend)
//	require.Equal(t, e.Description, createdEvent.Description)
//	require.Equal(t, e.Duration, createdEvent.Duration)
//	require.Equal(t, e.ID, createdEvent.ID)
//	require.Equal(t, e.Title, createdEvent.Title)
//}
//
//func testGetEventByID(t *testing.T, cfg *config.Config, createdEvent domain.Event) {
//	url := fmt.Sprintf("http://%s:%s/event/%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, createdEvent.ID)
//	b, statusCode := doRequest(t, http.MethodGet, url, []byte("{}"))
//	require.Equal(t, http.StatusOK, statusCode)
//
//	e := domain.Event{}
//	err := json.Unmarshal(b, &e)
//	require.Nil(t, err)
//
//	require.Equal(t, e.UserID, createdEvent.UserID)
//	require.Equal(t, e.IsNotificationSend, createdEvent.IsNotificationSend)
//	require.Equal(t, e.Description, createdEvent.Description)
//	require.Equal(t, e.Duration, createdEvent.Duration)
//	require.Equal(t, e.ID, createdEvent.ID)
//	require.Equal(t, e.Title, createdEvent.Title)
//}
//
//func doRequest(t *testing.T, method, url string, body []byte) ([]byte, int) {
//	req, err := http.NewRequest(method, url, bytes.NewReader(body))
//	if err != nil {
//		log.Fatal(err)
//	}
//	req.Header.Set("Content-Type", "application/json")
//	cl := http.Client{}
//	res, err := cl.Do(req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer res.Body.Close()
//
//	b, err := ioutil.ReadAll(res.Body)
//	require.Nil(t, err)
//
//	return b, res.StatusCode
//}
