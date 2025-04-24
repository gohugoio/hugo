---
title: Linux
description: Install Hugo on Linux.
categories: []
keywords: []
weight: 20
---

## Editions

{{% include "/_common/installation/01-editions.md" %}}

Unless your specific deployment needs require the extended/deploy edition, we recommend the extended edition.

{{% include "/_common/installation/02-prerequisites.md" %}}

{{% include "/_common/installation/03-prebuilt-binaries.md" %}}

## Package managers

### Snap

[Snap] is a free and open-source package manager for Linux. Available for [most distributions], snap packages are simple to install and are automatically updated.

The Hugo snap package is [strictly confined]. Strictly confined snaps run in complete isolation, up to a minimal access level that's deemed always safe. The sites you create and build must be located within your home directory, or on removable media.

To install the extended edition of Hugo:

```sh
sudo snap install hugo
```

To control automatic updates:

```sh
# disable automatic updates
sudo snap refresh --hold hugo

# enable automatic updates
sudo snap refresh --unhold hugo
```

To control access to removable media:

```sh
# allow access
sudo snap connect hugo:removable-media

# revoke access
sudo snap disconnect hugo:removable-media
```

To control access to SSH keys:

```sh
# allow access
sudo snap connect hugo:ssh-keys

# revoke access
sudo snap disconnect hugo:ssh-keys
```

{{% include "/_common/installation/homebrew.md" %}}

## Repository packages

Most Linux distributions maintain a repository for commonly installed applications.

> [!note]
> The Hugo version available in package repositories varies based on Linux distribution and release, and in some cases will not be the [latest version].
>
> Use one of the other installation methods if your package repository does not provide the desired version.

### Alpine Linux

To install the extended edition of Hugo on [Alpine Linux]:

```sh
doas apk add --no-cache --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community hugo
```

### Arch Linux

Derivatives of the [Arch Linux] distribution of Linux include [EndeavourOS], [Garuda Linux], [Manjaro], and others. To install the extended edition of Hugo:

```sh
sudo pacman -S hugo
```

### Debian

Derivatives of the [Debian] distribution of Linux include [elementary OS], [KDE neon], [Linux Lite], [Linux Mint], [MX Linux], [Pop!_OS], [Ubuntu], [Zorin OS], and others. To install the extended edition of Hugo:

```sh
sudo apt install hugo
```

You can also download Debian packages from the [latest release] page.

### Exherbo

To install the extended edition of Hugo on [Exherbo]:

1. Add this line to /etc/paludis/options.conf:

    ```text
    www-apps/hugo extended
    ```

1. Install using the Paludis package manager:

    ```sh
    cave resolve -x repository/heirecka
    cave resolve -x hugo
    ```

### Fedora

Derivatives of the [Fedora] distribution of Linux include [CentOS], [Red Hat Enterprise Linux], and others. To install the extended edition of Hugo:

```sh
sudo dnf install hugo
```

### Gentoo

Derivatives of the [Gentoo] distribution of Linux include [Calculate Linux], [Funtoo], and others. To install the extended edition of Hugo:

1. Specify the `extended` [USE] flag in /etc/portage/package.use/hugo:

    ```text
    www-apps/hugo extended
    ```

1. Build using the Portage package manager:

    ```sh
    sudo emerge www-apps/hugo
    ```

### NixOS

The NixOS distribution of Linux includes Hugo in its package repository. To install the extended edition of Hugo:

```sh
nix-env -iA nixos.hugo
```

### openSUSE

Derivatives of the [openSUSE] distribution of Linux include [GeckoLinux], [Linux Karmada], and others. To install the extended edition of Hugo:

```sh
sudo zypper install hugo
```

### Solus

The [Solus] distribution of Linux includes Hugo in its package repository. To install the extended edition of Hugo:

```sh
sudo eopkg install hugo
```

### Void Linux

To install the extended edition of Hugo on [Void Linux]:

```sh
sudo xbps-install -S hugo
```

{{% include "/_common/installation/04-build-from-source.md" %}}

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

[Alpine Linux]: https://alpinelinux.org/
[Arch Linux]: https://archlinux.org/
[Calculate Linux]: https://www.calculate-linux.org/
[CentOS]: https://www.centos.org/
[Debian]: https://www.debian.org/
[elementary OS]: https://elementary.io/
[EndeavourOS]: https://endeavouros.com/
[Exherbo]: https://www.exherbolinux.org/
[Fedora]: https://getfedora.org/
[Funtoo]: https://www.funtoo.org/
[Garuda Linux]: https://garudalinux.org/
[GeckoLinux]: https://geckolinux.github.io/
[Gentoo]: https://www.gentoo.org/
[KDE neon]: https://neon.kde.org/
[latest version]: https://github.com/gohugoio/hugo/releases/latest
[Linux Karmada]: https://linuxkamarada.com/
[Linux Lite]: https://www.linuxliteos.com/
[Linux Mint]: https://linuxmint.com/
[Manjaro]: https://manjaro.org/
[most distributions]: https://snapcraft.io/docs/installing-snapd
[MX Linux]: https://mxlinux.org/
[openSUSE]: https://www.opensuse.org/
[Pop!_OS]: https://pop.system76.com/
[Red Hat Enterprise Linux]: https://www.redhat.com/
[Snap]: https://snapcraft.io/
[Solus]: https://getsol.us/
[strictly confined]: https://snapcraft.io/docs/snap-confinement
[Ubuntu]: https://ubuntu.com/
[USE]: https://packages.gentoo.org/packages/www-apps/hugo
[Void Linux]: https://voidlinux.org/
[Zorin OS]: https://zorin.com/os/
