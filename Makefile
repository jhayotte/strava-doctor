.PHONY: all build clean run test

all: build

build:
	go build -o bin/strava-doctor ./...

clean:
	rm -rf bin

run: build
	./bin/strava-doctor

test:
	go test ./...

