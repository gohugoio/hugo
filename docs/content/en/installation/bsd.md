---
title: BSD
linkTitle: BSD
description: Install Hugo on BSD derivatives.
categories: [installation]
menu:
  docs:
    parent: installation
    weight: 50
toc: true
weight: 50
---
{{% readfile file="/installation/common/01-editions.md" %}}

{{% readfile file="/installation/common/02-prerequisites.md" %}}

{{% readfile file="/installation/common/03-prebuilt-binaries.md" %}}

## Repository packages

Most BSD derivatives maintain a repository for commonly installed applications. Please note that these repositories may not contain the [latest release].

[latest release]: https://github.com/gohugoio/hugo/releases/latest

### DragonFly BSD

[DragonFly BSD] includes Hugo in its package repository. This will install the extended edition of Hugo:

```sh
sudo pkg install gohugo
```

[DragonFly BSD]: https://www.dragonflybsd.org/

### FreeBSD

[FreeBSD] includes Hugo in its package repository. This will install the extended edition of Hugo:

```sh
sudo pkg install gohugo
```

[FreeBSD]: https://www.freebsd.org/

### NetBSD

[NetBSD] includes Hugo in its package repository. This will install the extended edition of Hugo:

```sh
sudo pkgin install go-hugo
```

[NetBSD]: https://www.netbsd.org/

### OpenBSD

[OpenBSD] includes Hugo in its package repository. This will prompt you to select which edition of Hugo to install:

```sh
doas pkg_add hugo
```

[OpenBSD]: https://www.openbsd.org/

{{% readfile file="/installation/common/04-docker.md" %}}

{{% readfile file="/installation/common/05-build-from-source.md" %}}

## Comparison

||Prebuilt binaries|Repository packages|Docker|Build from source
:--|:--:|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|
Easy to upgrade?|:heavy_check_mark:|varies|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|varies|:heavy_check_mark:|:heavy_check_mark:
Automatic updates?|:x:|varies|:x: [^1]|:x:
Latest version available?|:heavy_check_mark:|varies|:heavy_check_mark:|:heavy_check_mark:

[^1]: Possible but requires advanced configuration.
