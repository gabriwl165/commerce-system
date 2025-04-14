protoc:
	protoc --proto_path=../proto  --go_out=./proto/ --go_opt=paths=source_relative --go-grpc_out=./proto/ --go-grpc_opt=paths=source_relative ingestor.proto

build:
	docker build -t ingestor-api . 

docker-run:
	docker run -p 50051:50051 ingestor-api