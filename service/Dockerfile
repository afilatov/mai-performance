FROM golang:latest

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o /app/service-entrypoint ./cmd/service

ENTRYPOINT ["/app/service-entrypoint"]
