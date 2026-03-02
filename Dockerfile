FROM golang:1.23.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/

RUN go build -o main ./internal/main.go

EXPOSE 8080

ENTRYPOINT ["/app/main"]
