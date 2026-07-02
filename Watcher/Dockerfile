FROM golang:1.26.2-alpine AS builder

WORKDIR /app
RUN apk add git

COPY go.* .
RUN go mod download

COPY *.go .
COPY internal internal

RUN go build -o watcher .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/watcher .
COPY config.toml .

EXPOSE 8080

CMD ["./watcher"]