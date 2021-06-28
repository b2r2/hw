package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/sql"

	// init database driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	log := logrus.New()
	events, err := sqlstorage.New(log, ctx, "host=localhost port=5432 user=calendar password=calendar dbname=postgres sslmode=disable")
	require.NoError(t, err)
	require.NotNil(t, events)

	event1 := storage.Event{
		Title:       "First",
		ID:          1,
		Description: "bla",
		UserID:      1,
		Start:       time.Date(2021, 0o6, 15, 22, 15, 21, 0, time.UTC),
		Stop:        time.Date(2021, 0o6, 16, 11, 31, 5, 0, time.UTC),
	}
	event2 := storage.Event{
		Title:       "Second",
		ID:          2,
		Description: "bla-bla",
		UserID:      2,
		Start:       time.Date(2021, 0o6, 15, 22, 15, 21, 0, time.UTC),
		Stop:        time.Date(2021, 0o6, 16, 11, 31, 5, 0, time.UTC),
	}
	event2Update := storage.Event{
		Title:       "Second Update",
		ID:          2,
		Description: "bla-bla update",
		UserID:      2,
		Start:       time.Date(2021, 0o6, 15, 22, 15, 21, 0, time.UTC),
		Stop:        time.Date(2021, 0o6, 16, 11, 31, 5, 0, time.UTC),
	}
	_ = event2Update

	id, err := events.Create(ctx, &event1)
	require.NoError(t, err)
	require.Equal(t, id, int32(1))

	id, err = events.Create(ctx, &event2)
	require.NoError(t, err)
	require.Equal(t, id, int32(2))

	allEvents, err := events.ListAll(ctx)
	require.NoError(t, err)
	require.Len(t, allEvents, 2)

	err = events.Update(ctx, int32(1), &event2Update)
	require.NoError(t, err)

	event, err := events.Get(ctx, int32(1))
	require.NoError(t, err)
	require.Equal(t, int32(1), event.ID)

	event, err = events.Get(ctx, int32(123))
	require.Nil(t, event)
	require.Error(t, err)

	id, err = events.Create(ctx, &storage.Event{})
	require.NoError(t, err)
	require.Equal(t, int32(3), id)

	err = events.DeleteAll(ctx)
	require.NoError(t, err)
}
