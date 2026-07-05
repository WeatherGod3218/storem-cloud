FROM u

FROM golang:1.26.2-alpine AS builder

WORKDIR /app
RUN apk add git

COPY go.* .
RUN go mod download

COPY *.go .
COPY internal internal

RUN go build -o server .

FROM alpine:latest

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]