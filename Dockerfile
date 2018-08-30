# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

FROM golang:1.10.3-alpine3.7 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/github.com/gohugoio/hugo
RUN apk add --no-cache \
    git \
    musl-dev && \
  go get github.com/golang/dep/cmd/dep
COPY . /go/src/github.com/gohugoio/hugo/
RUN dep ensure -vendor-only && \
  go install -ldflags '-s -w'

# ---

FROM scratch
COPY --from=build /go/bin/hugo /hugo
WORKDIR /site
VOLUME  /site
EXPOSE  1313
ENTRYPOINT [ "/hugo" ]
CMD [ "--help" ]
