FROM golang:latest

RUN mkdir /app
COPY ./ /app
WORKDIR /app

RUN ulimit -c unlimited
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go mod vendor

RUN go build -gcflags=all="-N -l" -o /app/coredumper debugging.go

ENTRYPOINT ["/bin/bash"]
