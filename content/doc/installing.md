{
    "title": "Installing Hugo",
    "Pubdate": "2013-07-01"
}
Hugo is written in GoLang with support for Windows, Linux, FreeBSD and OSX.

The latest release can be found at [hugo releases](https://github.com/spf13/hugo/releases).
We currently build for Windows, Linux, FreeBSD and OS X for x64
and 386 architectures. 

Installation is very easy. Simply download the appropriate version for your
platform. Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally you should install it somewhere in your path for easy use. `/usr/local/bin` 
is the most probable location.

*Hugo has no external dependencies.*

Installation is very easy. Simply download the appropriate version for your
platform. 

## Installing from source

Make sure you have a recent version of go installed. Hugo requires go 1.1+.

    git clone https://github.com/spf13/hugo
    cd hugo
    go build -o hugo main.go

