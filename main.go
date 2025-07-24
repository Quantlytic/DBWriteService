package main

import (
	"fmt"
	"os"

	"github.com/Quantlytic/DBWriteService/pkg/kafkaconsumer"
)

func main() {
	err := kafkaconsumer.RunConsumer()
	if err != nil {
		fmt.Printf("Error running consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Consumer finished successfully")
}
