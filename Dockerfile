FROM golang:alpine3.6

ENV GOPATH /go

RUN \
  adduser -h /site -s /sbin/nologin -u 1000 -D hugo && \
  apk add --no-cache dumb-init && \
  apk add --no-cache --virtual .build-deps \
    git \
    make && \
  go get github.com/kardianos/govendor && \
  govendor get github.com/gohugoio/hugo && \
  cd $GOPATH/src/github.com/gohugoio/hugo && \
  make install test && \
  rm -rf $GOPATH/src/* && \
  apk del .build-deps

USER hugo

WORKDIR /site

EXPOSE 1313

ENTRYPOINT ["/usr/bin/dumb-init", "--", "hugo"]

CMD [ "--help" ]

