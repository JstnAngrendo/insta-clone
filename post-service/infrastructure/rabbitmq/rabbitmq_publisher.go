package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Publisher struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewPublisher(amqpURL, queueName string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Publisher) Publish(body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("[Publisher] Error marshaling body: %v", err)
		return err
	}

	log.Printf("[Publisher] Publishing to queue %s: %s", p.queue.Name, data)

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		log.Printf("[Publisher] Error publishing: %v", err)
		return err
	}

	log.Printf("[Publisher] Message published successfully")
	return nil
}
