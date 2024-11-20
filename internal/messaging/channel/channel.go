package channel

import (
	"fmt"
	"log"
	"sync"

	"github.com/cryptkeeperhq/cryptkeeper/internal/messaging" // replace with the correct import path
)

type ChannelProducer struct {
	topics map[string]chan []byte
	mu     sync.RWMutex
}

func NewChannelProducer() *ChannelProducer {
	return &ChannelProducer{topics: make(map[string]chan []byte)}
}

func (p *ChannelProducer) SendMessage(topic string, message []byte) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if ch, ok := p.topics[topic]; ok {
		ch <- message
		return nil
	}
	return fmt.Errorf("topic %s not found", topic)
}

func (p *ChannelProducer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, ch := range p.topics {
		close(ch)
	}
	return nil
}

type ChannelConsumer struct {
	producer *ChannelProducer
	handlers map[string]messaging.MessageHandler
	wg       sync.WaitGroup
}

func NewChannelConsumer(producer *ChannelProducer, handlers map[string]messaging.MessageHandler) *ChannelConsumer {
	return &ChannelConsumer{producer: producer, handlers: handlers}
}

func (c *ChannelConsumer) Consume(topics []string) {
	for _, topic := range topics {
		ch := make(chan []byte, 100)
		c.producer.mu.Lock()
		c.producer.topics[topic] = ch
		c.producer.mu.Unlock()

		c.wg.Add(1)
		go func(topic string, ch chan []byte) {
			defer c.wg.Done()
			for message := range ch {
				if handler, exists := c.handlers[topic]; exists {
					if err := handler.HandleMessage(topic, message); err != nil {
						log.Printf("Failed to handle message for topic %s: %v", topic, err)
					}
				}
			}
		}(topic, ch)
	}
}

func (c *ChannelConsumer) Close() error {
	c.wg.Wait()
	return nil
}
