---
title: Linux
description: Install Hugo on Linux.
categories: [installation]
keywords: []
menu:
  docs:
    parent: installation
    weight: 30
weight: 30
toc: true
---
{{% include "installation/_common/01-editions.md" %}}

{{% include "installation/_common/02-prerequisites.md" %}}

{{% include "installation/_common/03-prebuilt-binaries.md" %}}

## Package managers

### Snap

[Snap] is a free and open-source package manager for Linux. Available for [most distributions], snap packages are simple to install and are automatically updated.

The Hugo snap package is [strictly confined]. Strictly confined snaps run in complete isolation, up to a minimal access level that’s deemed always safe. The sites you create and build must be located within your home directory, or on removable media.

This will install the extended edition of Hugo:

```sh
sudo snap install hugo
```

To enable or revoke access to removable media:

```sh
sudo snap connect hugo:removable-media
sudo snap disconnect hugo:removable-media
```

To enable or revoke access to SSH keys:

```sh
sudo snap connect hugo:ssh-keys
sudo snap disconnect hugo:ssh-keys
```

[most distributions]: https://snapcraft.io/docs/installing-snapd
[strictly confined]: https://snapcraft.io/docs/snap-confinement
[Snap]: https://snapcraft.io/

{{% include "installation/_common/homebrew.md" %}}

## Repository packages

Most Linux distributions maintain a repository for commonly installed applications.

{{% note %}}
Due to Long Term Release (LTR) guidelines, most Linux package repositories will not contain the [latest release].

[latest release]: https://github.com/gohugoio/hugo/releases/latest
{{% /note %}}

### Arch Linux

Derivatives of the [Arch Linux] distribution of Linux include [EndeavourOS], [Garuda Linux], [Manjaro], and others. This will install the extended edition of Hugo:

```sh
sudo pacman -S hugo
```

[Arch Linux]: https://archlinux.org/
[EndeavourOS]: https://endeavouros.com/
[Manjaro]: https://manjaro.org/
[Garuda Linux]: https://garudalinux.org/

### Debian

Derivatives of the [Debian] distribution of Linux include [elementary OS], [KDE neon], [Linux Lite], [Linux Mint], [MX Linux], [Pop!_OS], [Ubuntu], [Zorin OS], and others. This will install the extended edition of Hugo:

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

Derivatives of the [Fedora] distribution of Linux include [CentOS], [Red Hat Enterprise Linux], and others. This will install the extended edition of Hugo:

```sh
sudo dnf install hugo
```

[CentOS]: https://www.centos.org/
[Fedora]: https://getfedora.org/
[Red Hat Enterprise Linux]: https://www.redhat.com/

### Gentoo

Derivatives of the [Gentoo] distribution of Linux include [Calculate Linux], [Funtoo], and others. Follow the instructions below to install the extended edition of Hugo:

1. Specify the `sass` [USE] flag in /etc/portage/package.use/mysql:

    ```text
    www-apps/hugo sass
    ```

2. Build using the Portage package manager:

    ```sh
    sudo emerge www-apps/hugo
    ```

[Calculate Linux]: https://www.calculate-linux.org/
[Funtoo]: https://www.funtoo.org/
[Gentoo]: https://www.gentoo.org/
[USE]: https://wiki.gentoo.org/wiki/Hugo#USE_flags

### openSUSE

Derivatives of the [openSUSE] distribution of Linux include [GeckoLinux], [Linux Karmada], and others. This will install the extended edition of Hugo:

```sh
sudo zypper install hugo
```

[GeckoLinux]: https://geckolinux.github.io/
[Linux Karmada]: https://linuxkamarada.com/
[openSUSE]: https://www.opensuse.org/

### Solus

The [Solus] distribution of Linux includes Hugo in its package repository. This will install the _standard_ edition of Hugo:

```sh
sudo eopkg install hugo
```

[Solus]: https://getsol.us/

{{% include "installation/_common/04-build-from-source.md" %}}

## Comparison

||Prebuilt binaries|Package managers|Repository packages|Build from source
:--|:--:|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|varies|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^1]|varies|:heavy_check_mark:
Automatic updates?|:x:|varies [^2]|:x:|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|varies|:heavy_check_mark:

[^1]: Easy if a previous version is still installed.
[^2]: Snap packages are automatically updated. Homebrew requires advanced configuration.
