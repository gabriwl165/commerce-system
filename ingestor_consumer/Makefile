protoc:
	protoc --proto_path=../proto  --go_out=./proto/ --go_opt=paths=source_relative --go-grpc_out=./proto/ --go-grpc_opt=paths=source_relative ingestor.proto

docker:
	docker start commerce-system-kafdrop-1 commerce-system-kafka-1 commerce-system-zookeeper-1

build:
	docker build -t ingestor-consumer . 

docker-run:
	docker run ingestor-consumer 