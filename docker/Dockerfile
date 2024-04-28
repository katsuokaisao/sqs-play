FROM golang:1.20

RUN apt-get update && apt-get install -y \
    curl \
    vim \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY ./go/go.mod ./go/go.sum ./
RUN go mod download

COPY ./go .
