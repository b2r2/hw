package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	Base
	Events
}

type Base interface {
	Close(ctx context.Context) error
}

type Events interface {
	Create(ctx context.Context, event *Event) (int32, error)
	Update(ctx context.Context, id int32, change *Event) error
	Get(ctx context.Context, id int32) (*Event, error)
	Delete(ctx context.Context, id int32) error
	DeleteAll(ctx context.Context) error
	ListAll(ctx context.Context) ([]*Event, error)
	ListDay(ctx context.Context, date time.Time) ([]*Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]*Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]*Event, error)
	IsTimeBusy(ctx context.Context, start, stop time.Time, excludeID int32) (bool, error)
}

type Event struct {
	ID           int32
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int32
	Notification *time.Duration
}

var ErrNotExistsEvent = errors.New("no such event")
