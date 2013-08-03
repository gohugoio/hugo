---
title: "Installing Hugo"
Pubdate: "2013-07-01"
---

Hugo is written in GoLang with support for Windows, Linux, FreeBSD and OSX.

The latest release can be found at [hugo releases](https://github.com/spf13/hugo/releases).
We currently build for Windows, Linux, FreeBSD and OS X for x64
and 386 architectures. 

## Installing Hugo (binary)

Installation is very easy. Simply download the appropriate version for your
platform from [hugo releases](https://github.com/spf13/hugo/releases).
Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally you should install it somewhere in your path for easy use. `/usr/local/bin` 
is the most probable location.

*the Hugo executible has no external dependencies.*

## Installing from source

### Dependencies

* Git
* Go 1.1+
* Mercurial
* Bazaar

### Get directly from Github:

    go get github.com/spf13/hugo

### Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

## Contributing

Please see the [contributing guide](/doc/contributing)
