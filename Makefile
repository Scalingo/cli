all:
	go build -o appsdeck appsdeck

get:
	go get appsdeck

fmt:
	go fmt ./...
