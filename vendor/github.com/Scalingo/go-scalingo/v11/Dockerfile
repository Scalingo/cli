FROM golang:1.26.8

MAINTAINER Étienne Michon "etienne@scalingo.com"

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/Scalingo/go-scalingo

CMD $GOPATH/bin/go-scalingo
