package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// Message represents a message received from Redis pub/sub
type Message struct {
	Channel string
	Pattern string
	Payload string
}

// Subscription represents a Redis pub/sub subscription
type Subscription struct {
	pubsub interface{} // Can be *redis.PubSub for standalone or cluster
}

// Close closes the subscription
func (s *Subscription) Close() error {
	if s.pubsub != nil {
		if ps, ok := s.pubsub.(interface{ Close() error }); ok {
			return ps.Close()
		}
	}
	return nil
}

// Channel returns a go channel for receiving messages
func (s *Subscription) Channel() <-chan *Message {
	if s.pubsub == nil {
		return nil
	}

	// Type assertion for standalone client
	if ps, ok := s.pubsub.(interface{ Channel() <-chan *redis.Message }); ok {
		ch := ps.Channel()
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

	return nil
}

// ReceiveMessage receives a message from the subscription
func (s *Subscription) ReceiveMessage() (*Message, error) {
	if s.pubsub == nil {
		return nil, nil
	}

	// Type assertion for standalone client
	if ps, ok := s.pubsub.(interface {
		ReceiveMessage(context.Context) (*redis.Message, error)
	}); ok {
		msg, err := ps.ReceiveMessage(context.Background())
		if err != nil {
			return nil, err
		}
		return &Message{
			Channel: msg.Channel,
			Pattern: msg.Pattern,
			Payload: msg.Payload,
		}, nil
	}

	return nil, nil
}

// Ping sends a PING to the server
func (s *Subscription) Ping() error {
	if s.pubsub == nil {
		return nil
	}

	// Type assertion for standalone client
	if ps, ok := s.pubsub.(interface{ Ping(context.Context) error }); ok {
		return ps.Ping(context.Background())
	}

	return nil
}
