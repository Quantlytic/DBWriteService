package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Quantlytic/DBWriteService/internal/config"
	"github.com/Quantlytic/DBWriteService/pkg/kafkaconsumer"

	"github.com/Quantlytic/DBWriteService/pkg/minioclient"
	"github.com/Quantlytic/DBWriteService/pkg/models"
)

func main() {
	appConfig := config.Load()

	// Create Kafka Consumer
	cfg := kafkaconsumer.KafkaConsumerConfig{
		BootstrapServers: appConfig.KafkaBrokers,
		GroupID:          "minio-write-service",
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

	onMessage := func(bucketName string) kafkaconsumer.OnMessage {
		return func(msg string, commit func()) {
			fmt.Printf("Received message: %s\n", msg)

			var apiData models.ApiData
			err := json.Unmarshal([]byte(msg), &apiData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to unmarshal message: %v\n", err)
				return
			}

			objName := fmt.Sprintf("%s/%s/%s", time.Now().Format("2006-01-02"), apiData.Exchange, apiData.Symbol)
			exists, err := minioClient.ObjectExists(context.Background(), bucketName, objName)

			if err != nil {
				log.Fatalf("Failed to check if object exists: %v\n", err)
				return
			}

			if !exists {
				minioClient.WriteBytesToFile(context.Background(), bucketName, objName, []byte(apiData.Data))
			} else {
				minioClient.AppendBytesToFile(context.Background(), bucketName, objName, []byte(apiData.Data))
			}
		}
	}

	topics := []kafkaconsumer.TopicConfig{
		{
			TopicName: "stock-quotes",
			Callback:  onMessage("stock-quotes"),
		},
		{
			TopicName: "stock-trades",
			Callback:  onMessage("stock-trades"),
		},
	}

	consumer.SubscribeTopics(topics)
}
