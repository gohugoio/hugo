
---
aliases:
- /doc/installing/
lastmod: 2016-01-04
date: 2013-07-01
menu:
  main:
    parent: getting started
next: /overview/usage
prev: /overview/quickstart
title: Installing Hugo
weight: 20
---

Hugo is written in [Go][] with support for multiple platforms.

The latest release can be found at [Hugo Releases](https://github.com/spf13/hugo/releases).
We currently provide pre-built binaries for
<i class="fa fa-windows"></i>&nbsp;Windows,
<i class="fa fa-linux"></i>&nbsp;Linux,
<i class="fa freebsd-19px"></i>&nbsp;FreeBSD
and <i class="fa fa-apple"></i>&nbsp;OS&nbsp;X (Darwin)
for x64, i386 and ARM architectures.

Hugo may also be compiled from source wherever the Go compiler tool chain can run, e.g. for other operating systems including DragonFly BSD, OpenBSD, Plan&nbsp;9 and Solaris.  See http://golang.org/doc/install/source for the full set of supported combinations of target operating systems and compilation architectures.

## Installing Hugo (binary)

Installation is very easy. Simply download the appropriate version for your
platform from [Hugo Releases](https://github.com/spf13/hugo/releases).
Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally, you should install it somewhere in your `PATH` for easy use.
`/usr/local/bin` is the most probable location.

On macOS, if you have [Homebrew](http://brew.sh/), installation is even
easier: just run `brew update && brew install hugo`.

For a more detailed explanation follow the corresponding installation guides:

- [Installation on macOS]({{< relref "tutorials/installing-on-mac.md" >}})
- [Installation on Windows]({{< relref "tutorials/installing-on-windows.md" >}})

### Installing Pygments (optional)

The Hugo executable has one *optional* external dependency for source code highlighting (Pygments).

If you want to have source code highlighting using the [highlight shortcode](/extras/highlighting/),
you need to install the Python-based Pygments program. The procedure is outlined on the [Pygments home page](http://pygments.org/).

## Upgrading Hugo

Upgrading Hugo is as easy as downloading and replacing the executable you’ve
placed in your `PATH`.

## Installing Hugo on Linux from native packages

### Arch Linux

You can install Hugo from the [Arch user repository](https://aur.archlinux.org/) on Arch Linux or derivatives such as Manjaro.

    sudo pacman -S yaourt
    yaourt -S hugo

Be aware that Hugo is built from source. This means that additional tools like [Git](https://git-scm.com/) and [Go](https://golang.org/doc/install) will be installed as well.

### Debian and Ubuntu

Hugo has been included in Debian and Ubuntu since 2016, and thus installing Hugo is as simple as:

    sudo apt install hugo

Pros:

* Native Debian/Ubuntu package maintained by Debian Developers
* Pre-installed bash completion script and man pages for best interactive experience

Cons:

* Might not be the latest version, especially if you are using an older stable version (e.g., Ubuntu 16.04&nbsp;LTS).
  Until backports and PPA are available, you may consider installing the Hugo snap package to get the latest version of Hugo, as described below.

### Fedora and Red Hat

* https://copr.fedorainfracloud.org/coprs/spf13/Hugo/ (updated to Hugo v0.16)
* https://copr.fedorainfracloud.org/coprs/daftaupe/hugo/ (updated to Hugo v0.19)

See also [this discussion](https://discuss.gohugo.io/t/solved-fedora-copr-repository-out-of-service/2491).

### Snap package for Hugo

In any of the [Linux distributions that support snaps](http://snapcraft.io/docs/core/install):

    snap install hugo

> Note: Hugo-as-a-snap can write only inside the user’s `$HOME` directory—and gvfs-mounted directories owned by the user—because of Snaps’ confinement and security model.
> More information is also available [in this related GitHub issue](https://github.com/spf13/hugo/issues/3143).

## Installing from source

### Prerequisite tools for downloading and building source code

* [Git](http://git-scm.com/)
* [Go][] 1.8+
* [govendor][]

### Vendored Dependencies

Hugo uses [govendor][] to vendor dependencies, but we don't commit the vendored packages themselves to the Hugo git repository.
Therefore, a simple `go get` is not supported since `go get` is not vendor-aware.
You **must use govendor** to fetch Hugo's dependencies.

### Fetch from GitHub

    go get github.com/kardianos/govendor
    govendor get github.com/spf13/hugo

`govendor get` will fetch Hugo and all its dependent libraries to
`$HOME/go/src/github.com/spf13/hugo`, and compile everything into a final `hugo`
(or `hugo.exe`) executable, which you will find sitting inside
`$HOME/go/bin/`, all ready to go!

*Windows users: where you see the `$HOME` environment variable above, replace it with `%USERPROFILE%`.*


*Note: For syntax highlighting using the [highlight shortcode](/extras/highlighting/),
you need to install the Python-based [Pygments](http://pygments.org/) program.*

## Contributing

Please see the [contributing guide](/doc/contributing/) if you are interested in
working with the Hugo source or contributing to the project in any way.

[Go]: http://golang.org/
[govendor]: https://github.com/kardianos/govendor
