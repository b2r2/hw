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
	t.Skip()
	log := logrus.New()
	events := sqlstorage.New(log)
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
		UserID:      1,
		Start:       time.Date(2021, 0o6, 15, 22, 15, 21, 0, time.UTC),
		Stop:        time.Date(2021, 0o6, 16, 11, 31, 5, 0, time.UTC),
	}
	event2Update := storage.Event{
		Title:       "Second Update",
		ID:          2,
		Description: "bla-bla update",
		UserID:      1,
		Start:       time.Date(2021, 0o6, 15, 22, 15, 21, 0, time.UTC),
		Stop:        time.Date(2021, 0o6, 16, 11, 31, 5, 0, time.UTC),
	}

	ctx := context.Background()

	id, err := events.Create(ctx, event1)
	require.NoError(t, err)
	require.Equal(t, id, 1)

	id, err = events.Create(ctx, event2)
	require.NoError(t, err)
	require.Equal(t, id, uint64(1))

	allEvents, err := events.ListAll(ctx)
	require.NoError(t, err)
	require.Len(t, allEvents, 2)

	err = events.Update(ctx, 1, event2Update)
	require.NoError(t, err)

	event, err := events.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, event.ID, 1)

	event, err = events.Get(ctx, 123)
	require.Nil(t, event)
	require.ErrorIs(t, storage.ErrNotExistsEvent, err)

	id, err = events.Create(ctx, storage.Event{})
	require.NoError(t, err)
	require.Equal(t, id, 2)

	err = events.DeleteAll(ctx)
	require.NoError(t, err)
}
