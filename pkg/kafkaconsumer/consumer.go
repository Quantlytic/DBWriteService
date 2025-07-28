package kafkaconsumer

import (
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	MIN_COMMIT_COUNT = 10
	POLLING_RATE     = 50 // TODO: Optimize later
)

// Called to process messages. Commit func signals the message is processed and ready to commit to kafka
type OnMessage func(msg string, commit func())

type KafkaConsumerConfig struct {
	BootstrapServers string
	GroupID          string
	AutoOffsetReset  string
}

type KafkaConsumer struct {
	consumer           *kafka.Consumer
	cfg                KafkaConsumerConfig
	topicCallbacks     map[string]OnMessage
	partitionsToCommit []kafka.TopicPartition
	commitMutex        sync.Mutex
}

type TopicConfig struct {
	TopicName string
	Callback  OnMessage
}

// createConsumer initializes a new Kafka consumer with the provided configuration
func NewKafkaConsumer(cfg KafkaConsumerConfig) (*KafkaConsumer, error) {
	// Create confluent kafka-go consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.BootstrapServers,
		"group.id":           cfg.GroupID,
		"auto.offset.reset":  cfg.AutoOffsetReset,
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:       consumer,
		cfg:            cfg,
		topicCallbacks: make(map[string]OnMessage),
	}, nil
}

func (c *KafkaConsumer) Close() {
	if c.consumer != nil {
		c.consumer.Close()
	}
}

// Commit stores the offsets to be committed
func (c *KafkaConsumer) Commit() error {
	c.commitMutex.Lock()
	defer c.commitMutex.Unlock()

	if len(c.partitionsToCommit) == 0 {
		return nil
	}

	offsetsToCommit := make([]kafka.TopicPartition, len(c.partitionsToCommit))
	copy(offsetsToCommit, c.partitionsToCommit)

	// Reset the slice
	c.partitionsToCommit = c.partitionsToCommit[:0]

	for i := range offsetsToCommit {
		offsetsToCommit[i].Offset++
	}

	_, err := c.consumer.CommitOffsets(offsetsToCommit)
	if err != nil {
		return fmt.Errorf("failed to commit offsets: %w", err)
	}

	return nil
}

// Given a slice of topic names + callbacks, subscribe to those topics with kafka
func (c *KafkaConsumer) SubscribeTopics(topics []TopicConfig) error {
	// Pre-Allocate slice of topic names
	topicNames := make([]string, 0, len(topics))

	// Store callbacks for each topic
	for _, topic := range topics {
		if topic.Callback == nil {
			return fmt.Errorf("callback for topic %s is nil", topic.TopicName)
		}

		// Store mapping of topic name to callback
		c.topicCallbacks[topic.TopicName] = topic.Callback

		// Append topic name to the slice
		topicNames = append(topicNames, topic.TopicName)
	}

	// Subscribe to topics in kafka
	err := c.consumer.SubscribeTopics(topicNames, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	run := true

	for run {
		ev := c.consumer.Poll(POLLING_RATE)
		switch e := ev.(type) {
		case *kafka.Message:
			// Call the callback for the topic
			if callback, ok := c.topicCallbacks[*e.TopicPartition.Topic]; ok {
				commitFunc := func() {
					c.commitMutex.Lock()
					defer c.commitMutex.Unlock()
					c.partitionsToCommit = append(c.partitionsToCommit, e.TopicPartition)
				}
				go callback(string(e.Value), commitFunc)
			} else {
				fmt.Printf("No callback found for topic %s\n", *e.TopicPartition.Topic)
			}

		case kafka.PartitionEOF:
			fmt.Printf("Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			run = false
		}
	}

	fmt.Printf("Subscribed to topics: %v\n", topicNames)
	return nil
}
