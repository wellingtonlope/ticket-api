test:
	go test ./...

server:
	go run ./framework/rest/server.go

build:
	go build ./framework/rest/server.go
