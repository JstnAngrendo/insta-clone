FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o timeline-service
EXPOSE 8082
CMD ["./timeline-service"]