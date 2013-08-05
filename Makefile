all:
	go build -o appsdeck-cli appsdeck/cli

fmt:
	go fmt appsdeck/cli
	go fmt appsdeck/auth
	go fmt appsdeck/constatns
	go fmt appsdeck/api
