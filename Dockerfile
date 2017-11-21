FROM golang:1.9.2-alpine3.6 AS build

RUN apk add --no-cache \
    git \
    musl-dev \
 && go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/gohugoio/hugo/
WORKDIR /go/src/github.com/gohugoio/hugo
RUN dep ensure -vendor-only
COPY . /go/src/github.com/gohugoio/hugo/
RUN go install -ldflags '-s -w'

FROM alpine:3.6

RUN adduser -h /site -s /sbin/nologin -u 1000 -D hugo \
 && apk add --no-cache \
    dumb-init
COPY --from=build /go/bin/hugo /bin/hugo
USER hugo
WORKDIR /site
VOLUME /site
EXPOSE 1313

ENTRYPOINT ["/usr/bin/dumb-init", "--", "hugo"]
CMD [ "--help" ]
