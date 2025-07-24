FROM golang:1.24.3 AS builder

WORKDIR /app
RUN apt-get update && apt-get install -y git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user-service .

EXPOSE 8080

CMD ["./user-service"]
