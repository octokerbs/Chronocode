FROM golang:1.23 AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/chronocode/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /main .

EXPOSE 8080

CMD ["./main"]