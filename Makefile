
# Adds build information from git repo
#
# as suggested by tatsushid in
# https://github.com/spf13/hugo/issues/540

COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/spf13/hugo/hugolib.CommitHash ${COMMIT_HASH} -X github.com/spf13/hugo/hugolib.BuildDate ${BUILD_DATE}"

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

