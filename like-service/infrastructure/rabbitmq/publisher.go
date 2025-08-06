package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewPublisher(amqpURL string, queueName string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &Publisher{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func NewPublisherFromConn(conn *amqp.Connection, queueName string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (p *Publisher) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("[Publisher] Failed to publish message: %v", err)
		return err
	}

	log.Printf("[Publisher] Published message to queue %s: %s", p.queueName, body)
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
