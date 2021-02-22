package rabbitmq

import (
	"github.com/streadway/amqp"
)

type consumer struct {
	topicName string
	connUrl   string
	conn      *amqp.Connection
	channel   *amqp.Channel
}

type Consumer interface {
	Connect() error
	Consume(f func(b []byte) error) error
}

func NewConsumer(url, topicName string) Consumer {
	return &consumer{
		topicName: topicName,
		connUrl:   url,
	}
}

func (m *consumer) Connect() error {
	conn, err := amqp.Dial(m.connUrl)
	if err != nil {
		return err
	}
	m.conn = conn

	m.channel, err = m.conn.Channel()
	if err != nil {
		return err
	}

	err = m.channel.ExchangeDeclare(m.topicName, "topic", true, false, false, false, nil)
	if err != nil {
		return nil
	}

	return nil
}

func (m *consumer) Consume(f func(b []byte) error) error {
	_, err := m.channel.QueueDeclare(QueueEvents, true, false, false, false, nil)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	err = m.channel.QueueBind(QueueEvents, "#", ExchangeEvents, false, nil)
	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	msgs, err := m.channel.Consume(m.topicName, "random-key", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		if err = f(msg.Body); err != nil {
			if err = msg.Ack(false); err != nil {
				return err
			}

			return err
		}

		if err = msg.Ack(true); err != nil {
			return err
		}
	}

	return nil
}
