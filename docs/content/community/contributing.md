---
aliases:
- /doc/contributing/
- /meta/contributing/
date: 2013-07-01
menu:
  main:
    parent: community
next: /tutorials/automated-deployments
prev: /community/press
title: Contributing to Hugo
weight: 30
---

All contributions to Hugo are welcome. Whether you want to scratch an itch, or simply contribute to the project, feel free to pick something from the roadmap
or contact [spf13](http://spf13.com/) about what may make sense
to do next.

You should fork the project and make your changes.  *We encourage pull requests to discuss code changes.*


When you're ready to create a pull request, be sure to:

  * Have test cases for the new code.  If you have questions about how to do it, please ask in your pull request.
  * Run `go fmt`
  * Squash your commits into a single commit.  `git rebase -i`.  It's okay to force update your pull request.
  * Make sure `go test ./...` passes, and `go build` completes.  Our [Travis CI loop](https://travis-ci.org/spf13/hugo) will catch most things that are missing.  The exception: Windows.  We run on Windows from time to time, but if you have access, please check on a Windows machine too.

## Contribution Overview

1. Fork Hugo from https://github.com/spf13/hugo
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Commit passing tests to validate changes.
5. Run `go fmt`
6. Squash commits into a single (or logically grouped) commits (`git rebase -i`)
7. Push to the branch (`git push origin my-new-feature`)
8. Create new Pull Request


# Building from source

## Clone locally (for contributors):

    git clone https://github.com/spf13/hugo
    cd hugo
    go get

Because Go expects all of your libraries to be found in either
`$GOROOT` or `$GOPATH`, it's helpful to symlink the project to one
of the following paths:

 * `ln -s /path/to/your/hugo $GOPATH/src/github.com/spf13/hugo`
 * `ln -s /path/to/your/hugo $GOROOT/src/pkg/github.com/spf13/hugo`

## Running Hugo

    cd /path/to/hugo
    go install github.com/spf13/hugo/hugo
    go run main.go

## Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

