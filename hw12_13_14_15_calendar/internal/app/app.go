package app

import (
	"context"
	"errors"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrNoUserID    = errors.New("no user id of the event")
	ErrEmptyTitle  = errors.New("no title of the event")
	ErrStartInPast = errors.New("start time of the event in the past")
	ErrDateBusy    = errors.New("this time is already occupied by another event")
)

type App interface {
	CreateEvent(ctx context.Context, event *storage.Event) (id int32, err error)
	UpdateEvent(ctx context.Context, id int32, change *storage.Event) error
	DeleteEvent(ctx context.Context, id int32) error
	DeleteAllEvent(ctx context.Context) error
	GetEvent(ctx context.Context, id int32) (*storage.Event, error)
	ListAllEvents(ctx context.Context) ([]*storage.Event, error)
	ListDayEvents(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListWeekEvents(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListMonthEvents(ctx context.Context, date time.Time) ([]*storage.Event, error)
}

type app struct {
	logger  logger.Logger
	storage storage.Storage
}

func New(logger logger.Logger, storage storage.Storage) App {
	return &app{
		logger,
		storage,
	}
}

func (a *app) CreateEvent(ctx context.Context, event *storage.Event) (id int32, err error) {
	if event.UserID == 0 {
		err = ErrNoUserID
		return
	}
	if event.Title == "" {
		err = ErrEmptyTitle
		return
	}
	if event.Start.After(event.Stop) {
		event.Start, event.Stop = event.Stop, event.Start
	}
	if time.Now().After(event.Start) {
		err = ErrStartInPast
		return
	}
	isBusy, err := a.storage.IsTimeBusy(ctx, event.Start, event.Stop, 0)
	if err != nil {
		return
	}
	if isBusy {
		err = ErrDateBusy
		return
	}

	return a.storage.Create(ctx, event)
}

func (a *app) UpdateEvent(ctx context.Context, id int32, change *storage.Event) error {
	if change.Title == "" {
		return ErrEmptyTitle
	}
	if change.Start.After(change.Stop) {
		change.Start, change.Stop = change.Stop, change.Start
	}
	if time.Now().After(change.Start) {
		return ErrStartInPast
	}
	isBusy, err := a.storage.IsTimeBusy(ctx, change.Start, change.Stop, id)
	if err != nil {
		return err
	}
	if isBusy {
		return ErrDateBusy
	}

	return a.storage.Update(ctx, id, change)
}

func (a *app) DeleteEvent(ctx context.Context, id int32) error {
	return a.storage.Delete(ctx, id)
}

func (a *app) DeleteAllEvent(ctx context.Context) error {
	return a.storage.DeleteAll(ctx)
}

func (a *app) GetEvent(ctx context.Context, id int32) (*storage.Event, error) {
	return a.storage.Get(ctx, id)
}

func (a *app) ListAllEvents(ctx context.Context) ([]*storage.Event, error) {
	return a.storage.ListAll(ctx)
}

func (a *app) ListDayEvents(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListDay(ctx, date)
}

func (a *app) ListWeekEvents(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListWeek(ctx, date)
}

func (a *app) ListMonthEvents(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListMonth(ctx, date)
}
