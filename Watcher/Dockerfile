FROM golang:1.26.2-alpine AS builder

WORKDIR /app
RUN apk add git

COPY go.* .
RUN go mod download
COPY internal internal

COPY *.go .

RUN go build -o watcher .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/watcher .
COPY config.yaml .
COPY credentials.yaml .

CMD ["./watcher"]