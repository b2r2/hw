package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	memorystorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/sirupsen/logrus"
)

func TestHandler(t *testing.T) {
	log := logrus.New()
	// db, err := sqlstorage.New(log, context.Background(), "host=localhost port=5432 user=calendar password=calendar dbname=calendar sslmode=disable")
	db := memorystorage.New(log)
	// require.NoError(t, err)
	require.NotNil(t, db)
	handler := NewHandler(log, app.New(log, db))
	router := NewRouter(log, handler, "test create")

	var result struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data,omitempty"`
		Error string `json:"error,omitempty"`
		Code  int    `json:"code"`
	}

	t.Run("create | OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		event := newEvent(1, time.Now().Add(time.Hour))
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(event)
		require.NoError(t, err)

		r := httptest.NewRequest("POST", "/api/v1/create", body)
		router.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err = json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, 1, result.Data.ID)
	})

	t.Run("error 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(nil)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/api/v1/create", body)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		event := newEvent(1, time.Now().Add(time.Hour))
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(event)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/api/v1/create", body)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("update | OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := new(bytes.Buffer)
		updateEvent := newEvent(1, time.Now().Add(time.Hour*3))
		err := json.NewEncoder(body).Encode(updateEvent)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/api/v1/update/1", body)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update | 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(nil)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/api/v1/update/1", body)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("update | 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		event := newEvent(12, time.Now().Add(time.Hour))
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(event)
		require.NoError(t, err)
		r := httptest.NewRequest("POST", "/api/v1/update/12", body)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("list all | OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listAll", nil)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		var res JSONResponse
		err := json.NewDecoder(w.Body).Decode(&res)
		require.NoError(t, err)
		require.Equal(t, len(res.Data.([]interface{})), 1)
	})

	t.Run("delete event | OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/delete/1", nil)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("delete event | 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/delete/12", nil)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("delete all | OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/deleteAll", nil)
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func newEvent(id int, t time.Time) storage.Event {
	return storage.Event{
		ID:           id,
		Title:        "event",
		Start:        t,
		Stop:         t.Add(time.Hour * 15),
		Description:  "some desc",
		UserID:       int32(id),
		Notification: nil,
	}
}
