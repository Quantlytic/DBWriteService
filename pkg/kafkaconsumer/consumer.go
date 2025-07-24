package kafkaconsumer

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const MIN_COMMIT_COUNT = 10

func CreateConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "10.99.97.145:9092",
		"group.id":          "db-write-service-group",
		"auto.offset.reset": "smallest"})

	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s", err)
	}

	return consumer, nil
}

func RunConsumer() error {
	c, err := CreateConsumer()
	if err != nil {
		fmt.Printf("Error initializing consumer: %s\n", err)
		return err
	}

	defer c.Close()

	// Subscribe to the topic
	fmt.Println("Subscribing to topic 'stock-data-raw'")

	if err := c.SubscribeTopics([]string{"stock-data-raw"}, nil); err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return fmt.Errorf("failed to subscribe to topic: %s", err)
	}

	fmt.Println("Subscribed to topic 'stock-data-raw'")

	run := true
	msg_count := 0

	for run {
		ev := c.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			msg_count += 1
			if msg_count%MIN_COMMIT_COUNT == 0 {
				// async commit every MIN_COMMIT_COUNT messages
				go func() {
					c.Commit()
				}()
			}
			fmt.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))

		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
			// default:
			// 	fmt.Printf("Ignored %v\n", e)
		}
	}

	return nil
}
