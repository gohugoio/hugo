# Developer: Maik Ellerbrock <opensource@frapsoft.com>
#
# GitHub:  https://github.com/ellerbrock
# Twitter: https://twitter.com/frapsoft
# Docker:  https://hub.docker.com/u/ellerbrock
# Quay:    https://quay.io/user/ellerbrock

FROM golang:alpine3.6

LABEL maintainer "Maik Ellerbrock <opensource@frapsoft.com>"

ENV GOPATH /go

RUN \
  adduser -h /site -s /sbin/nologin -u 1000 -D hugo && \
  apk add --no-cache \
    dumb-init \
    git && \
  go get github.com/gohugoio/hugo && \
  cd $GOPATH/src/github.com/gohugoio/hugo && \
  go install && \
  cd $GOPATH && \
  rm -rf pkg src && \
  apk del --no-cache git go

USER    hugo
WORKDIR /site
VOLUME  /site
EXPOSE  1313

ENTRYPOINT ["/usr/bin/dumb-init", "--", "hugo"]
CMD [ "--help" ]

