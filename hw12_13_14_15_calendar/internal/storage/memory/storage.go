package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
)

type store struct {
	log     logger.Logger
	mu      sync.Mutex
	counter int32
	data    map[int32]storage.Event
}

func (s *store) Close(_ context.Context) error {
	return nil
}

func New(log logger.Logger) storage.Storage {
	return &store{
		data: make(map[int32]storage.Event),
		log:  log,
	}
}

func (s *store) Create(_ context.Context, event *storage.Event) (int32, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	event.ID = id
	s.data[id] = storage.Event{
		ID:           id,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
	s.log.Traceln("create new event:", id)
	return id, nil
}

func (s *store) Update(_ context.Context, id int32, change *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.data[id]
	if !ok {
		return storage.ErrNotExistsEvent
	}

	event.Title = change.Title
	event.Start = change.Start
	event.Stop = change.Stop
	event.Description = change.Description
	event.Notification = change.Notification
	s.data[id] = event
	s.log.Traceln("update event:", id)

	return nil
}

func (s *store) Delete(_ context.Context, id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; ok {
		delete(s.data, id)
	} else {
		return storage.ErrNotExistsEvent
	}
	s.log.Traceln("deleted event:\n", id)
	return nil
}

func (s *store) DeleteAll(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[int32]storage.Event)

	s.log.Traceln("deleted all events")
	return nil
}

func (s *store) Get(_ context.Context, id int32) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, ok := s.data[id]; ok {
		return &event, nil
	}
	return nil, storage.ErrNotExistsEvent
}

func (s *store) ListAll(_ context.Context) ([]*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*storage.Event, 0, len(s.data))
	for _, event := range s.data {
		result = append(result, &event)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) ListDay(_ context.Context, date time.Time) ([]*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*storage.Event, 0, len(s.data))
	year, month, day := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, eventDay := event.Start.Date()
		if eventYear == year && eventMonth == month && eventDay == day {
			result = append(result, &event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) ListWeek(_ context.Context, date time.Time) ([]*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*storage.Event, 0, len(s.data))
	year, week := date.ISOWeek()
	for _, event := range s.data {
		eventYear, eventWeek := event.Start.ISOWeek()
		if eventYear == year && eventWeek == week {
			result = append(result, &event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) ListMonth(_ context.Context, date time.Time) ([]*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []*storage.Event
	year, month, _ := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, _ := event.Start.Date()
		if eventYear == year && eventMonth == month {
			result = append(result, &event)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result, nil
}

func (s *store) IsTimeBusy(_ context.Context, start, stop time.Time, excludeID int32) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.data {
		if event.ID != excludeID && event.Start.Before(stop) && event.Stop.After(start) {
			return true, nil
		}
	}
	return false, nil
}

func (s *store) newID() int32 {
	s.counter++
	return s.counter
}
