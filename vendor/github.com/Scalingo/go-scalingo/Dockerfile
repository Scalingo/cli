FROM golang:1.11
MAINTAINER Étienne Michon "etienne@scalingo.com"

RUN go get github.com/cespare/reflex

WORKDIR $GOPATH/src/github.com/Scalingo/go-scalingo

CMD $GOPATH/bin/go-scalingo
