build:
	cd ingestor_api/ && make build
	cd ingestor_consumer/ && make build

compose:
	docker compose up -d