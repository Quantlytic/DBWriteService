package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/Quantlytic/DBWriteService/internal/config"
	"github.com/Quantlytic/DBWriteService/pkg/kafkaconsumer"
)

func main() {
	appConfig := config.Load()

	cfg := kafkaconsumer.KafkaConsumerConfig{
		BootstrapServers: appConfig.KafkaBrokers,
		GroupID:          "db-write-service",
		AutoOffsetReset:  "earliest",
	}

	consumer, err := kafkaconsumer.NewKafkaConsumer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	fmt.Println("Kafka consumer created successfully")

	var (
		msgCount int
		mutex    sync.Mutex
	)

	messageProcessor := func(msg string, commit func()) {
		fmt.Printf("Received message: %s\n", msg)
		commit() // Stage this message for commit

		mutex.Lock()
		defer mutex.Unlock()
		msgCount++

		if msgCount >= 3 {
			fmt.Println("Processed 3 messages, committing now...")
			err := consumer.Commit()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to commit: %v\n", err)
			} else {
				fmt.Println("Commit successful.")
			}
			msgCount = 0 // Reset counter
		}
	}

	topics := []kafkaconsumer.TopicConfig{
		{
			TopicName: appConfig.KafkaTopic,
			Callback:  messageProcessor,
		},
	}

	consumer.SubscribeTopics(topics)
}
