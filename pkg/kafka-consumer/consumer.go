package kafka-consumer

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        // User-specific properties that you must set
        "bootstrap.servers": "<BOOTSTRAP SERVERS>",

        // Fixed properties
        "group.id":          "kafka-go-getting-started",
        "auto.offset.reset": "earliest"})

    if err != nil {
        fmt.Printf("Failed to create consumer: %s", err)
        os.Exit(1)
    }

    topic := "purchases"
    err = c.SubscribeTopics([]string{topic}, nil)
    // Set up a channel for handling Ctrl-C, etc
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    // Process messages
    run := true
    for run {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev, err := c.ReadMessage(100 * time.Millisecond)
            if err != nil {
                // Errors are informational and automatically handled by the consumer
                continue
            }
            fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
                *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
        }
    }

    c.Close()

}