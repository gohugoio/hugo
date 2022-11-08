---
title: Linux
linkTitle: Linux
description: Install Hugo on Linux.
categories: [installation]
menu:
  docs:
    parent: installation
    weight: 30
toc: true
weight: 30
---
{{% readfile file="/installation/common/01-flavors.md" %}}

{{% readfile file="/installation/common/02-prerequisites.md" %}}

{{% readfile file="/installation/common/03-prebuilt-binaries.md" %}}

## Package managers

### Snap

[Snap] is a free and open source package manager for Linux. Available for [most distributions], Snap packages are simple to install and are automatically updated. This will install the extended flavor of Hugo:

```sh
sudo snap install hugo
```

[most distributions]: https://snapcraft.io/docs/installing-snapd
[Snap]: https://snapcraft.io/

{{% readfile file="/installation/common/homebrew.md" %}}

## Repository packages

Most Linux distributions maintain a repository for commonly installed applications. Please note that these repositories may not contain the [latest release].

[latest release]: https://github.com/gohugoio/hugo/releases/latest

### Arch Linux

Derivatives of the [Arch Linux] distribution of Linux include [EndeavourOS], [Garuda Linux], [Manjaro], and others. This will install the extended flavor of Hugo:

```sh
sudo pacman -S hugo
```

[Arch Linux]: https://archlinux.org/
[EndeavourOS]: https://endeavouros.com/
[Manjaro]: https://manjaro.org/
[Garuda Linux]: https://garudalinux.org/

### Debian

Derivatives of the [Debian] distribution of Linux include [elementary OS], [KDE neon], [Linux Lite], [Linux Mint], [MX Linux], [Pop!_OS], [Ubuntu], [Zorin OS], and others. This will install the extended flavor of Hugo:

```sh
sudo apt install hugo
```

You can also download Debian packages from the [latest release] page.

[Debian]: https://www.debian.org/
[elementary OS]: https://elementary.io/
[KDE neon]: https://neon.kde.org/
[Linux Lite]: https://www.linuxliteos.com/
[Linux Mint]: https://linuxmint.com/
[MX Linux]: https://mxlinux.org/
[Pop!_OS]: https://pop.system76.com/
[Ubuntu]: https://ubuntu.com/
[Zorin OS]: https://zorin.com/os/

### Fedora

Derivatives of the [Fedora] distribution of Linux include [CentOS], [Red Hat Enterprise Linux], and others. This will install the extended flavor of Hugo:


```sh
sudo dnf install hugo
```

[CentOS]: https://www.centos.org/
[Fedora]: https://getfedora.org/
[Red Hat Enterprise Linux]: https://www.redhat.com/

### openSUSE

Derivatives of the [openSUSE] distribution of Linux include [GeckoLinux], [Linux Karmada], and others. This will install the extended flavor of Hugo:


```sh
sudo zypper install hugo
```

[GeckoLinux]: https://geckolinux.github.io/
[Linux Karmada]: https://linuxkamarada.com/
[openSUSE]: https://www.opensuse.org/

### Solus

The [Solus] distribution of Linux includes Hugo in its package repository. This will install the _standard_ flavor of Hugo:

```sh
sudo eopkg install hugo
```

[Solus]: https://getsol.us/home/

{{% readfile file="/installation/common/04-docker.md" %}}

{{% readfile file="/installation/common/05-build-from-source.md" %}}

## Comparison

||Prebuilt binaries|Package managers|Repository packages|Docker|Build from source
:--|:--:|:--:|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|varies|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^1]|varies|:heavy_check_mark:|:heavy_check_mark:
Automatic updates?|:x:|varies [^2]|:x:|:x: [^3]|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|varies|:heavy_check_mark:|:heavy_check_mark:

[^1]: Easy if a previous version is still installed.
[^2]: Snap packages are automatically updated. Homebrew requires advanced configuration.
[^3]: Possible but requires advanced configuration.
