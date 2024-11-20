package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/cryptkeeperhq/cryptkeeper/internal/messaging" // replace with the correct import path
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil
}

func (p *KafkaProducer) SendMessage(topic string, message []byte) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	})
	return err
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	handlers map[string]messaging.MessageHandler
}

func NewKafkaConsumer(brokers []string, handlers map[string]messaging.MessageHandler) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer, handlers: handlers}, nil
}

func (c *KafkaConsumer) Consume(topics []string) {
	for _, topic := range topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			log.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
		}
		for _, partition := range partitions {
			pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Failed to consume partition %d: %v", partition, err)
				continue
			}
			go func(pc sarama.PartitionConsumer) {
				for msg := range pc.Messages() {
					handler, exists := c.handlers[msg.Topic]
					if exists {
						if err := handler.HandleMessage(msg.Topic, msg.Value); err != nil {
							log.Printf("Error handling message from topic %s: %v", msg.Topic, err)
						}
					}
				}
			}(pc)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}
