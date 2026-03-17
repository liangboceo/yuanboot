package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// PubSub provides pub/sub operations following the pattern of KV, List, etc.
type PubSub struct {
	ops Ops
}

// Publish posts a message to the channel
func (ps PubSub) Publish(channel string, message interface{}) (int64, error) {
	return ps.ops.Publish(channel, message)
}

// Subscribe subscribes to the specified channels
func (ps PubSub) Subscribe(channels ...string) (*Subscription, error) {
	return ps.ops.Subscribe(channels...)
}

// PSubscribe subscribes to channels matching the given patterns
func (ps PubSub) PSubscribe(patterns ...string) (*Subscription, error) {
	return ps.ops.PSubscribe(patterns...)
}

// ReceiveMessage receives a single message from a channel
func (ps PubSub) ReceiveMessage(channel string) (*Message, error) {
	sub, err := ps.Subscribe(channel)
	if err != nil {
		return nil, err
	}
	defer func(sub *Subscription) {
		_ = sub.Close()
	}(sub)
	return sub.ReceiveMessage()
}

// ReceiveMessageTimeout receives a single message from a channel with timeout
func (ps PubSub) ReceiveMessageTimeout(channel string, timeout time.Duration) (*Message, error) {
	sub, err := ps.Subscribe(channel)
	if err != nil {
		return nil, err
	}
	defer func(sub *Subscription) {
		_ = sub.Close()
	}(sub)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Receive with timeout
	psSub := sub.pubsub
	if psSub == nil {
		return nil, nil
	}

	// Type assertion and receive
	if receiveMsg, ok := psSub.(interface {
		ReceiveMessage(context.Context) (*redis.Message, error)
	}); ok {
		msg, err := receiveMsg.ReceiveMessage(timeoutCtx)
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

// ReceiveMessages continuously receives messages from a channel
func (ps PubSub) ReceiveMessages(channel string, messageChan chan<- *Message) (*Subscription, error) {
	sub, err := ps.Subscribe(channel)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(messageChan)
		defer func(sub *Subscription) {
			_ = sub.Close()
		}(sub)

		psSub := sub.pubsub
		if psSub == nil {
			return
		}

		// Get channel
		if ch, ok := psSub.(interface{ Channel() <-chan *redis.Message }); ok {
			redisChan := ch.Channel()
			for msg := range redisChan {
				messageChan <- &Message{
					Channel: msg.Channel,
					Pattern: msg.Pattern,
					Payload: msg.Payload,
				}
			}
		}
	}()

	return sub, nil
}

// ReceiveMessagesWithPattern continuously receives messages matching a pattern
func (ps PubSub) ReceiveMessagesWithPattern(pattern string, messageChan chan<- *Message) (*Subscription, error) {
	sub, err := ps.PSubscribe(pattern)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(messageChan)
		defer func(sub *Subscription) {
			_ = sub.Close()
		}(sub)

		psSub := sub.pubsub
		if psSub == nil {
			return
		}

		// Get channel
		if ch, ok := psSub.(interface{ Channel() <-chan *redis.Message }); ok {
			redisChan := ch.Channel()
			for msg := range redisChan {
				messageChan <- &Message{
					Channel: msg.Channel,
					Pattern: msg.Pattern,
					Payload: msg.Payload,
				}
			}
		}
	}()

	return sub, nil
}

// SubscribeWithCallback subscribes to channels and invokes callback for each message
func (ps PubSub) SubscribeWithCallback(callback func(*Message), channels ...string) (*Subscription, error) {
	sub, err := ps.Subscribe(channels...)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func(sub *Subscription) {
			_ = sub.Close()
		}(sub)

		psSub := sub.pubsub
		if psSub == nil {
			return
		}

		// Get channel
		if ch, ok := psSub.(interface{ Channel() <-chan *redis.Message }); ok {
			redisChan := ch.Channel()
			for msg := range redisChan {
				callback(&Message{
					Channel: msg.Channel,
					Pattern: msg.Pattern,
					Payload: msg.Payload,
				})
			}
		}
	}()

	return sub, nil
}

// PSubscribeWithCallback subscribes to patterns and invokes callback for each message
func (ps PubSub) PSubscribeWithCallback(callback func(*Message), patterns ...string) (*Subscription, error) {
	sub, err := ps.PSubscribe(patterns...)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func(sub *Subscription) {
			_ = sub.Close()
		}(sub)

		psSub := sub.pubsub
		if psSub == nil {
			return
		}

		// Get channel
		if ch, ok := psSub.(interface{ Channel() <-chan *redis.Message }); ok {
			redisChan := ch.Channel()
			for msg := range redisChan {
				callback(&Message{
					Channel: msg.Channel,
					Pattern: msg.Pattern,
					Payload: msg.Payload,
				})
			}
		}
	}()

	return sub, nil
}
