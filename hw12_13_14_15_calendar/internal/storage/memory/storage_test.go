package memory_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("CRUD", func(t *testing.T) {
		log := logrus.New()
		events := memorystorage.New(log)
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
		id, err := events.Create(ctx, &event1)
		require.NoError(t, err)
		require.Equal(t, id, int32(1))

		id, err = events.Create(ctx, &event2)
		require.NoError(t, err)
		require.Equal(t, id, int32(2))

		allEvents, _ := events.ListAll(ctx)
		require.Len(t, allEvents, 2)

		err = events.Update(ctx, 2, &event2Update)
		require.NoError(t, err)

		allEvents, err = events.ListAll(ctx)
		require.NoError(t, err)
		require.Len(t, allEvents, 2)

		event, err := events.Get(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, event.ID, int32(1))
		require.Equal(t, event.Title, "First")

		_, err = events.Get(ctx, 5)
		require.ErrorIs(t, storage.ErrNotExistsEvent, err)

		err = events.Delete(ctx, 1)
		require.NoError(t, err)

		id, err = events.Create(ctx, &storage.Event{})
		require.NoError(t, err)
		require.Equal(t, id, int32(3))
	})
	t.Run("concurrent", func(t *testing.T) {
		l := 53
		log := logrus.New()
		events := memorystorage.New(log)
		var wg sync.WaitGroup
		for i := 0; i < l; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := events.Create(context.Background(), &storage.Event{})
				require.NoError(t, err)
			}()
		}
		wg.Wait()
		allEvents, err := events.ListAll(context.Background())
		require.NoError(t, err)
		require.Len(t, allEvents, l)
	})
}
