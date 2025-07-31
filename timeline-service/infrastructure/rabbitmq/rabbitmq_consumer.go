package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/entities"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/usecases"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	usecase usecases.TimelineUseCase
}

func NewConsumer(amqpURL string, usecase usecases.TimelineUseCase) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		usecase: usecase,
	}, nil
}

func (c *Consumer) StartConsuming(queueName string) error {
	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			var evt entities.PostCreatedEvent
			if err := json.Unmarshal(msg.Body, &evt); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			err := c.usecase.ProcessNewPost(context.Background(), evt)
			if err != nil {
				log.Printf("Error processing new post event: %v", err)
			} else {
				log.Println("Successfully processed new post event.")
			}
		}
	}()

	log.Printf("Consumer started on queue: %s", queueName)
	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
