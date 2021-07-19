package rmq

import (
	"encoding/json"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type Producer struct {
	log  logger.Logger
	name string
	conn *amqp.Connection
	ttl  int
}

func NewProducer(log logger.Logger, dsn, name string, ttl int) (*Producer, error) {
	var p Producer
	p.log = log
	p.name = name
	p.ttl = ttl

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	p.conn = conn

	return &p, nil
}

func (p *Producer) Publish(message *storage.Notification) error {
	ch, err := p.conn.Channel()
	if err != nil {
		p.log.Errorln("cannot create channel", err)
		return err
	}

	_, err = ch.QueueDeclare(p.name, false, false, false, false, amqp.Table{"x-message-ttl": p.ttl})
	if err != nil {
		p.log.Errorln("cannot create queue", err)
		return err
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if err := ch.Publish("", p.name, false, false, amqp.Publishing{ContentType: "text/plain", Body: data}); err != nil {
		p.log.Errorln("cannot publish message", err)
		return err
	}

	return nil
}
