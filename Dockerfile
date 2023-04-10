FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

FROM alpine:latest AS production

WORKDIR /app

COPY --from=builder /app/main .
COPY config /app/config

EXPOSE 53/tcp 53/udp

CMD ["./main focus --proxy"]