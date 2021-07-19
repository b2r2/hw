package rmq

import (
	"context"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type Consumer struct {
	name string
	log  logger.Logger
	conn *amqp.Connection
}

func NewConsumer(log logger.Logger, dsn, name string) (*Consumer, error) {
	var c Consumer
	c.log = log
	c.name = name
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return &c, nil
}

func (c *Consumer) Consumer(ctx context.Context, queue string) (<-chan Message, error) {
	messages := make(chan Message)
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		if err := ch.Close(); err != nil {
			c.log.Errorln(err)
		}
	}()
	delivery, err := ch.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			close(messages)
			c.log.Infoln("message channel closed")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-delivery:
				if err := d.Ack(false); err != nil {
					c.log.Errorln("cannot delivery message")
				}
				message := Message{context.TODO(), d.Body}
				select {
				case <-ctx.Done():
					return
				case messages <- message:
				}

			}
		}
	}()
	return messages, nil
}
