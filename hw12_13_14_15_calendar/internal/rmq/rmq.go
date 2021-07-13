package rmq

import (
	"context"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
)

type Scheduler struct {
	log     logger.Logger
	storage storage.Storage
	ticker  *time.Ticker
}

func NewScheduler(log logger.Logger, storage storage.Storage, ticker *time.Ticker) *Scheduler {
	return &Scheduler{log, storage, ticker}
}

type Message struct {
	ctx  context.Context
	Data []byte
}

func (r *Scheduler) Start(ctx context.Context, dsn, name string, ttl int) error {
	p, err := NewProducer(r.log, dsn, name, ttl)
	if err != nil {
		r.log.Errorln("cannot to connect to producer", err)
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.ticker.C:
				r.notify(ctx, p)
			}
		}
	}()

	<-ctx.Done()

	return nil
}

func (r *Scheduler) notify(ctx context.Context, p *Producer) {
	events, err := r.storage.ListNotifyEvents(ctx)
	if err != nil {
		r.log.Errorln("cannot events for notify", err)
	}

	for _, e := range events {
		message := e.GetNotification()
		err := p.Publish(message)
		if err != nil {
			r.log.Errorln("cannot publish message", err)
		}
	}
}
