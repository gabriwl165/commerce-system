FROM golang:1.24.0
WORKDIR /app

# Copy the rest of the application code
COPY . .

RUN go build -o ingestor-consumer cmd/kafka/main.go

RUN apt-get update && apt-get install -y ca-certificates openssl

RUN openssl req -x509 -newkey rsa:4096 -keyout tls_key.pem -out tls_cert.pem -days 365 -nodes -subj '/CN=localhost'

EXPOSE 8080

# Run the binary
CMD ["./ingestor-consumer"]