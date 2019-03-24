# GitHub:       https://github.com/gohugoio
# Twitter:      https://twitter.com/gohugoio
# Website:      https://gohugo.io/

FROM golang:1.11-stretch AS build


WORKDIR /go/src/github.com/gohugoio/hugo
RUN apt-get install \
    git gcc g++ binutils
COPY . /go/src/github.com/gohugoio/hugo/
ENV GO111MODULE=on
RUN go get -d .

ARG CGO=0
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux

# default non-existent build tag so -tags always has an arg
ARG BUILD_TAGS="99notag"
RUN go install -ldflags '-w -extldflags "-static"' -tags ${BUILD_TAGS}

# ---

FROM scratch
COPY --from=build /go/bin/hugo /hugo
ARG  WORKDIR="/site"
WORKDIR ${WORKDIR}
VOLUME  ${WORKDIR}
EXPOSE  1313
ENTRYPOINT [ "/hugo" ]
CMD [ "--help" ]
