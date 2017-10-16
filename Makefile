# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

PACKAGE = github.com/gohugoio/hugo
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`
LDFLAGS = -ldflags "-X ${PACKAGE}/hugolib.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/hugolib.BuildDate=${BUILD_DATE}"
NOGI_LDFLAGS = -ldflags "-X ${PACKAGE}/hugolib.BuildDate=${BUILD_DATE}"

# allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
GOEXE ?= go

.PHONY: vendor docker check fmt lint test test-race vet test-cover-html help
.DEFAULT_GOAL := help

vendor: ## Install govendor and sync Hugo's vendored dependencies
	${GOEXE} get github.com/kardianos/govendor
	govendor sync ${PACKAGE}

hugo: vendor ## Build hugo binary
	${GOEXE} build ${LDFLAGS} ${PACKAGE}

hugo-race: vendor ## Build hugo binary with race detector enabled
	${GOEXE} build -race ${LDFLAGS} ${PACKAGE}

install: vendor ## Install hugo binary
	${GOEXE} install ${LDFLAGS} ${PACKAGE}

hugo-no-gitinfo: LDFLAGS = ${NOGI_LDFLAGS}
hugo-no-gitinfo: vendor hugo ## Build hugo without git info

docker: ## Build hugo Docker container
	docker build -t hugo .
	docker rm -f hugo-build || true
	docker run --name hugo-build hugo ls /go/bin
	docker cp hugo-build:/go/bin/hugo .
	docker rm hugo-build

govendor: vendor # Deprecated: use "vendor" target
get: vendor # Deprecated: use "vendor"
gitinfo: hugo # Deprecated: use "hugo" target
install-gitinfo: install # Deprecated: use "install" target
no-git-info: hugo-no-gitinfo # Deprecated: use "hugo-no-gitinfo" target

check: test-race test386 fmt vet ## Run tests and linters

test386: ## Run tests in 32-bit mode
	GOARCH=386 govendor test +local

test: ## Run tests
	govendor test +local

test-race: ## Run tests with race detector
	govendor test -race +local

fmt: ## Run gofmt linter
	@for d in `govendor list -no-status +local | sed 's/github.com.gohugoio.hugo/./'` ; do \
		if [ "`gofmt -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ improperly formatted go files" && echo && exit 1; \
		fi \
	done

lint: ## Run golint linter
	@for d in `govendor list -no-status +local | sed 's/github.com.gohugoio.hugo/./'` ; do \
		if [ "`golint $$d | tee /dev/stderr`" ]; then \
			echo "^ golint errors!" && echo && exit 1; \
		fi \
	done

vet: ## Run go vet linter
	@if [ "`govendor vet +local | tee /dev/stderr`" ]; then \
		echo "^ go vet errors!" && echo && exit 1; \
	fi

test-cover-html: PACKAGES = $(shell govendor list -no-status +local | sed 's/github.com.gohugoio.hugo/./')
test-cover-html: ## Generate test coverage report
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		govendor test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	${GOEXE} tool cover -html=coverage-all.out

check-vendor: ## Verify that vendored packages match git HEAD
	@git diff-index --quiet HEAD vendor/ || (echo "check-vendor target failed: vendored packages out of sync" && echo && git diff vendor/ && exit 1)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
