---
title: BSD
description: Install Hugo on BSD derivatives.
categories: []
keywords: []
weight: 40
---

## Editions

{{% include "/_common/installation/01-editions.md" %}}

Unless your specific deployment needs require the extended/deploy edition, we recommend the extended edition.

{{% include "/_common/installation/02-prerequisites.md" %}}

{{% include "/_common/installation/03-prebuilt-binaries.md" %}}

## Repository packages

Most BSD derivatives maintain a repository for commonly installed applications. Please note that these repositories may not contain the [latest release].

[latest release]: https://github.com/gohugoio/hugo/releases/latest

### DragonFly BSD

[DragonFly BSD] includes Hugo in its package repository. To install the extended edition of Hugo:

```sh
sudo pkg install gohugo
```

[DragonFly BSD]: https://www.dragonflybsd.org/

### FreeBSD

[FreeBSD] includes Hugo in its package repository. To install the extended edition of Hugo:

```sh
sudo pkg install gohugo
```

[FreeBSD]: https://www.freebsd.org/

### NetBSD

[NetBSD] includes Hugo in its package repository. To install the extended edition of Hugo:

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

{{% include "/_common/installation/04-build-from-source.md" %}}

## Comparison

&nbsp;|Prebuilt binaries|Repository packages|Build from source
:--|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to upgrade?|:heavy_check_mark:|varies|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|varies|:heavy_check_mark:
Automatic updates?|:x:|varies|:x:
Latest version available?|:heavy_check_mark:|varies|:heavy_check_mark:
