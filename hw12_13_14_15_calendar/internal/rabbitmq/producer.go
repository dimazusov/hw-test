package rabbitmq

import (
	"github.com/streadway/amqp"
)

type producer struct {
	exchangeName string
	connUrl      string
	conn         *amqp.Connection
	channel      *amqp.Channel
}

type Producer interface {
	Connect() error
	SendMessage(b []byte) error
}

func NewProducer(url, exchangeName string) Producer {
	return &producer{
		exchangeName: exchangeName,
		connUrl:      url,
	}
}

func (m *producer) Connect() error {
	conn, err := amqp.Dial(m.connUrl)
	if err != nil {
		return err
	}
	m.conn = conn

	m.channel, err = m.conn.Channel()
	if err != nil {
		return err
	}

	err = m.channel.ExchangeDeclare(m.exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return nil
	}

	return nil
}

func (m *producer) SendMessage(b []byte) error {
	msg := amqp.Publishing{
		Body: b,
	}

	return m.channel.Publish(m.exchangeName, "random-key", false, false, msg)
}
