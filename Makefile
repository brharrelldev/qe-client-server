

build-server:
	go build -o bin/qe-server server/*.go

build-client:
	go build -o bin/qe-client client/*.go

all: build-server build-client