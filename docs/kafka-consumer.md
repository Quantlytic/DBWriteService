# Kafka Consumer
The design of the KafkaConsumer was chosen to decouple interacting with the kafka server from the confluent kafka go library and overall encapsulate the communication logic. 

## Package Structure
- `type KafkaConsumerConfig struct` - Configures static parameters that shouldn't change throughout execution, like the bootstrap server and groupID.
- `type OnMessage func` - Type of the handler function provided to `KafkaConsumer` upon subscribing to a topic.
- `type TopicConfig struct` - the name of a topic and its handler function, passed to the subscribe method.
- `type KafkaConsumer struct` - contains the underlying confluent kafka go Consumer, the configuration, a mapping of topic names to callbacks, a slice of topics/offsets ready to be committed, and a mutex for synchronization. Encapsulates kafka communication logic.

## Topic Handlers
The developer passes a `TopicName` and a `Callback` when they try to subscribe to a topic. The `Callback` controls how the data is handled upon a new message event. When a topic is subscribed to, the name and callback are inserted into a map to be accessed when a message is received from kafka. 

## Moving the Offset
When a message is received, the `Callback` should handle the contents of the message data. Once that is done, it can call the `commit` function which tells the kafka consumer that the offset can be moved up to this message.

The callback function requires a `commit func()` argument. This argument accesses the `TopicPartition` metadata that the confluent kafka go library provides when a message is consumed and adds it to the `partitionsToCommit` array. Next, the developer calls `KafkaConsumer.Commit()` to commit the offsets stored to kafka. The effect of this design choice is to allow the developer to batch data before committing to kafka. In the event of a failure, any data in the batch being processed will still reside in kafka. 
