# DBWriteService
Write service will subscribe to raw data kafka topic and write data to the database for long term storage

## Tasks
### Project Task List
- [ ] Create github action to deploy to kubernetes cluster
    - [x] Action - Build Image & Store in registry
    - [ ] Action - Deploy new release to prod
- [ ] Create Go package for reading data from kafka topic
    - [x] Basic functionality for subscribing and reading from topic
    - [ ] Add callback function to apply to each message
    - [ ] Allow configuration of topic, kafka ip, group, etc
- [ ] Create Go package for writing to Minio Bucket

### Progress
**July 24th -** I decided to start with creating the deployment. I have two minikube nodes - one running on a "server" that acts as my prod deployment. The other node runs locally for development. Since the service is going to interact with resources on the kubernetes cluster (kafka, minio) it made sense to make sure the service could run on the cluster.

On a push to the main branch, an image is built and pushed to github's image registry. Then, I can run `kubectl apply -f deploy/base/deploy.yaml` to test my deployment. Later on, I'll create the action that deploys the service to the prod node. To avoid pushing every time I want to build an image, I can build the image in minikube instead.

I was able to have my service connect to the kafka server and listen to a topic. I tested this using a kafka producer and was able to see the sent messages in the logs. To increase resusability of this package, I'm going to make the RunConsumer func accept a callback function that is applied to each message as it is received. I'll also have to make sure aspects like the kafka ip, group, and other initialization parameters are configurable instead of hardcoded.


## Setup
Run the following command to setup credentials needed to pull the image from github's image repository. 

```
kubectl create secret docker-registry ghcr-creds \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_PAT>
```

To create an image locally without pushing to the repository, run 
```
minikube image build -t ghcr.io/quantlytic/dbwriteservice:latest -f deploy/dockerfile .
```