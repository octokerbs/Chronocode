FROM golang:1.23.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/chronocode/server.go

COPY migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/app/main"]