# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

FROM golang:1.11-alpine3.7 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/gohugoio/hugo
RUN apk add --no-cache \
    git \
    musl-dev
COPY . /go/src/github.com/gohugoio/hugo/
RUN go install -ldflags '-s -w'

# ---

FROM scratch
COPY --from=build /go/bin/hugo /hugo
WORKDIR /site
VOLUME  /site
EXPOSE  1313
ENTRYPOINT [ "/hugo" ]
CMD [ "--help" ]
