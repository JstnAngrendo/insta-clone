package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/entities"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/repositories"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/usecases"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	usecase usecases.TimelineUseCase
	repo    repositories.RedisRepository
}

func NewConsumer(amqpURL string, uc usecases.TimelineUseCase) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Consumer{conn: conn, channel: ch, usecase: uc}, nil
}

func (c *Consumer) StartConsuming(queueName string) error {
	_, err := c.channel.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	msgs, err := c.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	log.Printf("[Consumer] Listening on queue %s", queueName)
	go func() {
		for d := range msgs {
			log.Printf("[Consumer] Raw message: %s", d.Body)
			var evt entities.PostCreatedEvent
			if err := json.Unmarshal(d.Body, &evt); err != nil {
				log.Printf("[Consumer] JSON unmarshal error: %v", err)
				continue
			}
			log.Printf("[Consumer] Parsed event: %+v", evt)
			if err := c.usecase.ProcessNewPost(context.Background(), evt); err != nil {
				log.Printf("[Consumer] UseCase error: %v", err)
			} else {
				log.Printf("[Consumer] Processed event: post %d for user %d", evt.PostID, evt.AuthorID)
			}
		}
	}()
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
