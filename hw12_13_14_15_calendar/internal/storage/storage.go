package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	Create(ctx context.Context, event *Event) (int, error)
	Update(ctx context.Context, id int, change *Event) error
	Get(ctx context.Context, id int) (*Event, error)
	Delete(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) error
	ListAll(ctx context.Context) ([]*Event, error)
	ListDay(ctx context.Context, date time.Time) ([]*Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]*Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]*Event, error)
	IsTimeBusy(ctx context.Context, start, stop time.Time, excludeID int) (bool, error)
	ListNotifyEvents(ctx context.Context) ([]*Event, error)
}

type Event struct {
	ID               int
	Title            string
	Start            time.Time
	Stop             time.Time
	Description      string
	UserID           int32
	NotificationTime *time.Duration
}

func (e *Event) GetNotification() *Notification {
	return &Notification{
		ID:        e.ID,
		Title:     e.Title,
		EventTime: e.Start,
		Owner:     e.UserID,
	}
}

type Notification struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	EventTime time.Time `json:"event_time"`
	Owner     int32     `json:"owner"`
}

func (n *Notification) Encode() ([]byte, error) {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(n)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *Notification) String() string {
	return fmt.Sprintf("id: %d, title: %s, time: %s, onwer: %d", n.ID, n.Title, n.EventTime, n.Owner)
}

var ErrNotExistsEvent = errors.New("no such event")
