build:
	go build -o bin/goraft

run: build
	./bin/goraft

runfollower: build
	./bin/goraft --listenaddr :3000 --leaderaddr :4000

test: 
	@go test -v ./...