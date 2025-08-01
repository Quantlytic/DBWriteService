# DBWriteService
Write service will subscribe to raw data kafka topic and write data to the database for long term storage

## Tasks
### Project Task List
- [ ] Create github action to deploy to kubernetes cluster
    - [x] Action - Build Image & Store in registry
    - [ ] Action - Deploy new release to prod
- [x] Create Go package for reading data from kafka topic
    - [x] Basic functionality for subscribing and reading from topic
    - [x] Add callback function to apply to each message
    - [x] Allow configuration of topic, kafka ip, group, etc
- [x] Create Go package for writing to Minio Bucket

### Progress
I'm keeping notes on progress and design decisions in the [docs directory](/docs/progress-log.md).

## Setup
Run the following command to setup credentials needed to pull the image from github's image repository. 

```
kubectl create secret docker-registry ghcr-creds \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_PAT>
```

To create an image locally without pushing to the repository, run 
- Temporarily set docker command context to inside minikube container: `eval $(minikube -p minikube docker-env)`
- Build docker image: `docker image build -t ghcr.io/quantlytic/dbwriteservice:latest -f deploy/dockerfile .`

## Testing Locally
While developing the service, it's quicker to run locally for testing instead of building an image each time.

Run the following in bash to port forward Kafka and Minio API:
```
# Kafka
kubectl port-forward -n kafka pod/quantlytic-kafka-dual-role-0 9092:9092
# Minio
kubectl port-forward -n minio-dev pod/minio 9000:9000
```

Next, add the following line to your `/etc/hosts` file to avoid DNS resolution errors:
```
127.0.0.1 quantlytic-kafka-dual-role-0.quantlytic-kafka-kafka-brokers.kafka.svc
```

You can now run the service locally against the kafka server in minikube. To run a Kafka producer, use the following command:
```
kubectl exec -it quantlytic-kafka-dual-role-0 -n kafka -- \
  /opt/kafka/bin/kafka-console-producer.sh \
    --bootstrap-server <kafka-dualrole-ip>:9092 \
    --topic stock-data-raw
```

### Local Test Steps
Following these steps will allow you to test the integration of your write service, kafka, and minio. 

0. Setup port forwards and DNS resolution. Configure environment variables seen in `.env.example`. Create a bucket in Minio Web portal.
1. Run Producer script and write 3 messages to kafka
2. Notice each message is printed to the console. When the third message is received, a Commit occurs.
3. View the bucket in Minio Web. Notice that there are now 3 files.