

build-server:
	go build -o bin/qe-server server/main.go

build-client:
	go build -o bin/qe-client client/main.go

all: build-server build-client