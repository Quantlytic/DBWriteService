package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Quantlytic/DBWriteService/internal/config"
	"github.com/Quantlytic/DBWriteService/pkg/kafkaconsumer"

	"github.com/Quantlytic/DBWriteService/pkg/minioclient"
)

func main() {
	appConfig := config.Load()

	// Create Kafka Consumer
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

	// Create minio writer
	minioClient, err := minioclient.NewMinioClient(appConfig.MinioEndpoint, appConfig.MinioAccessKey, appConfig.MinioSecretKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Minio client: %v\n", err)
		os.Exit(1)
	}

	var (
		msgCount int
		mutex    sync.Mutex
	)

	messageProcessor := func(msg string, commit func()) {
		fmt.Printf("Received message: %s\n", msg)

		// Write message to Minio
		err := minioClient.WriteBytesToFile(context.Background(), "raw-stock-data", fmt.Sprintf("test/%d", msgCount), []byte(msg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to Minio: %v\n",
				err)
			return
		}

		// Stage this message for commit
		commit()

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
