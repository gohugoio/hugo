FROM golang:1.8-alpine

ENV GOPATH /go
ENV USER root

RUN apk update && apk add git make

# pre-install known dependencies before the source, so we don't redownload them whenever the source changes
RUN go get github.com/kardianos/govendor \
 && govendor get github.com/spf13/hugo

COPY . $GOPATH/src/github.com/spf13/hugo

RUN cd $GOPATH/src/github.com/spf13/hugo \
 	&& make install test

ENTRYPOINT "/bin/sh"
