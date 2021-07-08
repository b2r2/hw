package rmq

import (
	"context"
	"fmt"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type Client struct {
	log   logger.Logger
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

var retry int

func New(l logger.Logger, dsn string, ttl int) (*Client, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("error rmq connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to make rmq channel: %w", err)
	}

	q, err := ch.QueueDeclare("notification",
		false,
		false,
		false,
		false,
		amqp.Table{"x-message-ttl": ttl})
	if err != nil {
		return nil, fmt.Errorf("failed to make rmq queue: %w", err)
	}

	return &Client{
		log:   l,
		conn:  conn,
		ch:    ch,
		queue: q,
	}, nil
}

func (c *Client) Notify(events []*storage.Event) {
	for _, event := range events {
		msg, err := event.GetNotification().Encode()
		if err != nil {
			c.log.Errorln("failed to encoding message:", err)
			continue
		}

		for i := 0; i < retry; i++ {
			if err := c.ch.Publish("",
				c.queue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg,
				},
			); err != nil {
				c.log.Errorln("failed to publish message:", err)
				continue
			}
			c.log.Infoln("send notification on", event.ID, ":", event.Title)
			break
		}
	}
}

func (c *Client) Send(ctx context.Context, sender func([]byte)) error {
	msg, err := c.ch.Consume(c.queue.Name, "sender", true, false, false, false, nil)
	if err != nil {
		c.log.Errorln("failed to send message:", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-msg:
			sender(m.Body)
		}
	}
}

func (c *Client) Close() error {
	if err := c.ch.Close(); err != nil {
		c.log.Errorln("error closing channel rmq:", err)
		return err
	}

	if err := c.conn.Close(); err != nil {
		c.log.Errorln("error closing connection rmq:", err)
		return err
	}

	return nil
}
