build:
	go build -o bin/goraft

run: build
	./bin/goraft