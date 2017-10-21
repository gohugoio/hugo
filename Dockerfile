FROM golang:1.9.0-alpine3.6 AS build

RUN apk add --no-cache --virtual git musl-dev
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/gohugoio/hugo
RUN dep ensure
ADD . /go/src/github.com/gohugoio/hugo/
RUN go install -ldflags '-s -w'

FROM alpine:3.6
RUN \
  adduser -h /site -s /sbin/nologin -u 1000 -D hugo && \
  apk add --no-cache \
    dumb-init
COPY --from=build /go/bin/hugo /bin/hugo
USER    hugo
WORKDIR /site
VOLUME  /site
EXPOSE  1313

ENTRYPOINT ["/usr/bin/dumb-init", "--", "hugo"]
CMD [ "--help" ]
