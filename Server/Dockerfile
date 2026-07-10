FROM node:22-alpine AS npmbuilder

WORKDIR /app

COPY web/package*.json ./

RUN npm ci

COPY web/ .

RUN npm run build

FROM golang:1.26.2-alpine AS gobuilder

WORKDIR /app
RUN apk add git

COPY go.* .
RUN go mod download
COPY internal internal
COPY docs docs

COPY *.go .

COPY --from=npmbuilder /app/dist/ ./web/dist/
RUN go build -o server .

FROM alpine:latest

WORKDIR /app
COPY --from=gobuilder /app/server .

EXPOSE 8080

CMD ["./server"]

