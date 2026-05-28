package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// Message represents a message received from Redis pub/sub
type Message struct {
	Channel string
	Pattern string
	Payload string
}

// Subscription represents a Redis pub/sub subscription
type Subscription struct {
	pubsub *redis.PubSub
}

// Close closes the subscription
func (s *Subscription) Close() error {
	if s.pubsub != nil {
		return s.pubsub.Close()
	}
	return nil
}

// Channel returns a go channel for receiving messages
func (s *Subscription) Channel() <-chan *Message {
	if s.pubsub == nil {
		return nil
	}

	ch := s.pubsub.Channel()
	msgChan := make(chan *Message, 100)
	go func() {
		defer close(msgChan)
		for msg := range ch {
			msgChan <- &Message{
				Channel: msg.Channel,
				Pattern: msg.Pattern,
				Payload: msg.Payload,
			}
		}
	}()
	return msgChan
}

// ReceiveMessage receives a message from the subscription
func (s *Subscription) ReceiveMessage() (*Message, error) {
	if s.pubsub == nil {
		return nil, nil
	}

	msg, err := s.pubsub.ReceiveMessage(context.Background())
	if err != nil {
		return nil, err
	}
	return &Message{
		Channel: msg.Channel,
		Pattern: msg.Pattern,
		Payload: msg.Payload,
	}, nil
}

// Ping sends a PING to the server
func (s *Subscription) Ping() error {
	if s.pubsub == nil {
		return nil
	}
	return s.pubsub.Ping(context.Background())
}
