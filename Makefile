.PHONY: all build docker-build up down test clean

all: build

build:
    go build -o operator main.go

docker-build:
    docker build -t local-operator:latest .

up: docker-build
    docker-compose up -d

down:
    docker-compose down

test:
    go test ./... -cover

clean:
    docker-compose down
    docker rmi local-operator:latest || true
    rm -f operator
