all:
	go build -o appsdeck-cli appsdeck/cli

get:
	go get github.com/Appsdeck/cli
	go get code.google.com/p/gopass

fmt:
	go fmt appsdeck/cli
	go fmt appsdeck/cli/auth
	go fmt appsdeck/cli/constants
	go fmt appsdeck/cli/api
	go fmt appsdeck/cli/apps
	go fmt appsdeck/cli/appdetect
	go fmt appsdeck/cli/cmd
