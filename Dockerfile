# GitHub:       https://github.com/gohugoio/
# DockerHub:    https://hub.docker.com/r/gohugoio/
# Quay:         https://quay.io/user/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

FROM golang:1.10-rc-alpine3.7 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/github.com/gohugoio/hugo

RUN \
  apk add --no-cache \
    git \
    musl-dev && \
  go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/gohugoio/hugo/

RUN dep ensure -vendor-only

COPY . /go/src/github.com/gohugoio/hugo/

RUN go install -ldflags '-s -w'

# ---

FROM scratch

COPY --from=build /go/bin/hugo /hugo

EXPOSE  1313

ENTRYPOINT [ "/hugo" ]
CMD [ "--help" ]

