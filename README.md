# About

This is a microservice application for a movie API. It consists of three services, the public facing Movie service, which calls the metadata service and the rating service. The microservices communicate through gRPC. 

# How to run

### Requirements
- [Docker](https://www.docker.com/)
- Make, while not compulsory, would make running this easier. 
- [grpcurl](https://github.com/fullstorydev/grpcurl), or an alternative tool which can send and recieve grpc requests

## Setup
#### 1. Run the following commands
---
- `make consul` - This runs the consul service registry
- `make jaeger` - This runs the jaeger tracing service

If the make tool is not available, go to `./Makefile` and copy the appropriate commands into the terminal.

#### 2. On 3 separate terminals, run the following commands from root directory
---
   1. `go run ./metadata/cmd/*.go`
   2. `go run ./rating/cmd/*.go`
   3. `go run ./movie/cmd/*.go`

For each service, you should recieve a log message which states the service has been started successfully.

## Sending Requests
#### 3. Use grpcurl or an alternative tool to begin sending requests.    
---

##### 1. Add metadata to a movie:

`grpcurl -plaintext -d '{"metadata":{"id":"1", "title": "Movie", "description":"This is a movie","director":"The Director"}}' localhost:8081 MetadataService/PutMetadata`

##### 2. Add a rating to the movie
`grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie", "user_id": "Aditya", "rating_value": 5}' localhost:8082 RatingService/PutRating`

##### 2(a). Add another rating to the movie - optional
`grpcurl -plaintext -d '{"record_id":"1", "record_type": "movie", "user_id": "Aditya", "rating_value": 7}' localhost:8082 RatingService/PutRating`

##### 3. Retrieve the movie details
`grpcurl -plaintext -d '{"movie_id":"1"}' localhost:8083 MovieService/GetMovieDetails`

##### 4. View the tracing on jaeger
Go to [localhost:16686](http://localhost:16686) to view the request trace using jaeger. 

# Features

- To view the proto file, check the [movie.proto](./api/movie.proto) file in `./api`. This will provide more details on the schemas used. 

- Currently the application uses an in memory db, however the code for using a MYSQL db is also implemented and working. The schema is visible in the [schema](./schema/schema.sql) file and the make command to setup the mysql docker container is also provided. Theoretically, simply changing the import package from memory to mysql in the main functions for metadata and rating should allow it to work, however this is untested and further modifications in the application may be required to use MySQL

- This runs on a consul service registry. To start a new instance of any service on a different port, run the `go run` command above with a `--port <PORT>` flag. (Make sure the port is not already in use!)

- This provides tracing of the request using jaeger. You can view the requests on [localhost:16686](http://localhost:16686)

# Todo 
- [ ] Fix prometheus, does not work for unknown reasons
- [ ] Provide screenshots of jaeger tracing and appropriate explanation in the readme
- [ ] Convert the updated code to docker images, and alter the setup to preferably use docker compose or something similar.
- [ ] Test files

# Learnings
- While developing this application, I used and learnt about the following
  - Kubernetes
  - Service Registries like Consul 
  - Tracing using Jaeger
  - Prometheus
  - Docker
  - gRPC
  - Structured logging tools such as zap
  - Kafka
  - Testing with mocks and mock generators such as mock from uber-go
  - Microservices architecture - controller service repository pattern.

# Links
[Book Link - Microservices With Go](https://www.packtpub.com/product/microservices-with-go/9781804617007)

Overall a well explained book, however there are instances where the book skips over some code changes and you need to sync up with the github repository before you can understand the next portion of the book.