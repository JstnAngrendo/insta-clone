package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/repositories"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PostLikeEvent struct {
	PostID uint   `json:"post_id"`
	Action string `json:"action"`
}

type Consumer struct {
	channel *amqp.Channel
	repo    repositories.PostRepository
}

func NewConsumer(conn *amqp.Connection, repo repositories.PostRepository) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{channel: ch, repo: repo}, nil
}

func (c *Consumer) Start(queueName string) {
	q, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Consumer queue declare: %v", err)
	}
	msgs, err := c.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Consumer register: %v", err)
	}
	go func() {
		for msg := range msgs {
			var evt PostLikeEvent
			if err := json.Unmarshal(msg.Body, &evt); err != nil {
				log.Printf("Consumer unmarshal: %v", err)
				continue
			}
			var delta int64
			if evt.Action == "like" {
				delta = 1
			} else if evt.Action == "unlike" {
				delta = -1
			} else {
				continue
			}
			if err := c.repo.UpdateLikeCount(context.Background(), evt.PostID, delta); err != nil {
				log.Printf("UpdateLikeCount error: %v", err)
			} else {
				log.Printf("like_count updated for post %d by %d", evt.PostID, delta)
			}
		}
	}()
	log.Printf("[Consumer] Listening on %s", queueName)
}
