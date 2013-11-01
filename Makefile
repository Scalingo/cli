all:
	go build -o appsdeck appsdeck/cli

get:
	go get github.com/Appsdeck/cli
	go get code.google.com/p/gopass

fmt:
	go fmt ./...
