FROM golang:1.24
MAINTAINER Étienne Michon "etienne@scalingo.com"

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/Scalingo/go-scalingo

CMD $GOPATH/bin/go-scalingo
