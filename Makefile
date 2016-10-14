
# Adds build information from git repo
#
# as suggested by tatsushid in
# https://github.com/spf13/hugo/issues/540

COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/spf13/hugo/hugolib.CommitHash=${COMMIT_HASH} -X github.com/spf13/hugo/hugolib.BuildDate=${BUILD_DATE}"

all: gitinfo

install: install-gitinfo

help:
	echo ${COMMIT_HASH}
	echo ${BUILD_DATE}

gitinfo:
	go build ${LDFLAGS} -o hugo main.go

install-gitinfo:
	go install ${LDFLAGS} ./...

no-git-info:
	go build -o hugo main.go

docker:
	docker build -t hugo .
	docker rm -f hugo-build || true
	docker run --name hugo-build hugo ls /go/bin
	docker cp hugo-build:/go/bin/hugo .
	docker rm hugo-build

govendor:
	go get -u github.com/kardianos/govendor
	go install github.com/kardianos/govendor
	govendor sync github.com/spf13/hugo

check: fmt vet test test-race

cyclo:
	@for d in `govendor list -no-status +local | sed 's/github.com.spf13.hugo/./'` ; do \
		if [ "`gocyclo -over 20 $$d | tee /dev/stderr`" ]; then \
			echo "^ cyclomatic complexity exceeds 20, refactor the code!" && echo && exit 1; \
		fi \
	done

fmt:
	@for d in `govendor list -no-status +local | sed 's/github.com.spf13.hugo/./'` ; do \
		if [ "`gofmt -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ improperly formatted go files" && echo && exit 1; \
		fi \
	done

lint:
	@for d in `govendor list -no-status +local | sed 's/github.com.spf13.hugo/./'` ; do \
		if [ "`golint $$d | tee /dev/stderr`" ]; then \
			echo "^ golint errors!" && echo && exit 1; \
		fi \
	done

get:
	go get -v -t ./...

test:
	govendor test +local

test-race:
	govendor test -race +local

vet:
	@if [ "`govendor vet +local | tee /dev/stderr`" ]; then \
		echo "^ go vet errors!" && echo && exit 1; \
	fi

