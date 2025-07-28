# Project Progress / Notes
**July 28th -** After learning how to create a consumer for my kafka topic, I realized two potential problems that I would run into in this project:
- If a consumer needed to subscribe to multiple topics, they would create multiple confluent kafka go consumers to connect to the same server
- If I wanted to batch raw stock data before writing to a file / minio, I will need to manually control committing. 
I changed the structure of the kafka-consumer interface and implementation to approach solving these problems. Read more [here!](/docs/kafka-consumer.md)

I'm still deciding what API(s) I'll use to gather stock data. Part of the solution will be using [this dataset](https://huggingface.co/datasets/bwzheng2010/yahoo-finance-data) to bootstrap my historical data. Going forward, I'd like to start collecting more granular data (possibly down to the minute) on top of end-of-day values. 

The next challenges I'll need to overcome:
- Defining what json data from kafka will look like
- Batching live data into parquet files for efficient storage


**July 24th -** I decided to start with creating the deployment. I have two minikube nodes - one running on a "server" that acts as my prod deployment. The other node runs locally for development. Since the service is going to interact with resources on the kubernetes cluster (kafka, minio) it made sense to make sure the service could run on the cluster.

On a push to the main branch, an image is built and pushed to github's image registry. Then, I can run `kubectl apply -f deploy/base/deploy.yaml` to test my deployment. Later on, I'll create the action that deploys the service to the prod node. To avoid pushing every time I want to build an image, I can build the image in minikube instead.

I was able to have my service connect to the kafka server and listen to a topic. I tested this using a kafka producer and was able to see the sent messages in the logs. To increase resusability of this package, I'm going to make the RunConsumer func accept a callback function that is applied to each message as it is received. I'll also have to make sure aspects like the kafka ip, group, and other initialization parameters are configurable instead of hardcoded.
