.PHONY: help consul jaeger db/create db/schema build/all docker/image

## help: Display this help message
help:

## consul: Runs docker container for consul service registry
consul: 
	@docker run --rm -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0

## jaeger: Runs docker container for Jaeger 
jaeger:
	@docker run -d --rm --name jaeger \
 	-e COLLECTOR_OTLP_ENABLED=true \
 	-p 6831:6831/udp \
 	-p 6832:6832/udp \
 	-p 5778:5778 \
 	-p 16686:16686 \
 	-p 4317:4317 \
 	-p 4318:4318 \
 	-p 14250:14250 \
 	-p 14268:14268 \
 	-p 14269:14269 \
 	-p 9411:9411 \
 	jaegertracing/all-in-one:1.37;

## db/create: Creates and runs docker container for mysql db
db/create:
	@docker run --name movieexample_db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=movieexample -p 3306:3306 -d mysql:latest

## db/schema: Migrates the schema into mysql db
db/schema:
	@docker exec -i movieexample_db mysql movieexample -h 0.0.0.0 -P 3306 --protocol=tcp -uroot -ppassword < schema/schema.sql

SERVICES=metadata rating movie

## builds go executables for metadata, rating, movie service
build/all:
	@for service in $(SERVICES); do \
		cd $$service && pwd && \
		GOOS=linux go build -o main cmd/main.go cmd/config.go && \
		cd ..; \
	done

## builds discord images for metadata, rating, movie service
docker/image:
	@for service in $(SERVICES); do \
		cd $$service && pwd && \
		docker build -t $$service . && \
		cd ..; \
	done