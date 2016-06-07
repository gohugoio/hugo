FROM golang:1.6.2
MAINTAINER Sven Dowideit <SvenDowideit@home.org.au>

ENV GOPATH /go
ENV USER root

# pre-install known dependencies before the source, so we don't redownload them whenever the source changes
RUN go get github.com/stretchr/testify/assert \
	&& go get github.com/kyokomi/emoji \
	&& go get github.com/bep/inflect \
	&& go get github.com/BurntSushi/toml \
	&& go get github.com/PuerkitoBio/purell \
	&& go get github.com/opennota/urlesc \
	&& go get github.com/dchest/cssmin \
	&& go get github.com/eknkc/amber \
	&& go get github.com/gorilla/websocket \
	&& go get github.com/kardianos/osext \
	&& go get github.com/miekg/mmark \
	&& go get github.com/mitchellh/mapstructure \
	&& go get github.com/russross/blackfriday \
	&& go get github.com/shurcooL/sanitized_anchor_name \
	&& go get github.com/spf13/afero \
	&& go get github.com/spf13/cast \
	&& go get github.com/spf13/jwalterweatherman \
	&& go get github.com/spf13/cobra \
	&& go get github.com/cpuguy83/go-md2man \
	&& go get github.com/inconshreveable/mousetrap \
	&& go get github.com/spf13/pflag \
	&& go get github.com/spf13/fsync \
	&& go get github.com/spf13/viper \
	&& go get github.com/kr/pretty \
	&& go get github.com/kr/text \
	&& go get github.com/magiconair/properties \
	&& go get golang.org/x/text/transform \
	&& go get golang.org/x/text/unicode/norm \
	&& go get github.com/yosssi/ace \
	&& go get github.com/spf13/nitro \
	&& go get github.com/fsnotify/fsnotify

COPY . /go/src/github.com/spf13/hugo
RUN go get -d -v github.com/spf13/hugo
RUN go install github.com/spf13/hugo

WORKDIR /go/src/github.com/spf13/hugo
RUN go get -d -v
RUN go build -o hugo main.go
RUN go test github.com/spf13/hugo/...

