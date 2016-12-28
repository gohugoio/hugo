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
	&& go get github.com/PuerkitoBio/urlesc \
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
	&& go get github.com/fortytw2/leaktest \
	&& go get github.com/fsnotify/fsnotify \
	&& go get github.com/bep/gitmap \
	&& go get github.com/nicksnyder/go-i18n/i18n

COPY . /go/src/github.com/spf13/hugo

RUN cd /go/src/github.com/spf13/hugo \
	&& go get -d -v \
	&& go install \
	&& go test github.com/spf13/hugo/...

