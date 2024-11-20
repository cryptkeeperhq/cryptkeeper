// Messaging interfaces for producer and consumer
package messaging

type MessageHandler interface {
	HandleMessage(topic string, message []byte) error
}

type Producer interface {
	SendMessage(topic string, message []byte) error
	Close() error
}

type Consumer interface {
	Consume(topics []string)
	Close() error
}
