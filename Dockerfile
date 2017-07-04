FROM golang:alpine3.6

ENV GOPATH /go

RUN \
  adduser -h /site -s /sbin/nologin -u 1000 -D hugo && \
  apk add --no-cache \
    dumb-init \
    git && \
  go get github.com/kardianos/govendor && \
  govendor get github.com/gohugoio/hugo && \
  cd $GOPATH/src/github.com/gohugoio/hugo && \
  go install && \
  cd $GOPATH && \
  rm -rf pkg src bin/govendor && \
  apk del --no-cache git go

USER    hugo
WORKDIR /site
VOLUME  /site
EXPOSE  1313

ENTRYPOINT ["/usr/bin/dumb-init", "--", "hugo"]
CMD [ "--help" ]

