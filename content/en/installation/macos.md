---
title: macOS
linkTitle: macOS
description: Install Hugo on macOS.
categories: [installation]
menu:
  docs:
    parent: installation
    weight: 20
toc: true
weight: 20
---
{{% readfile file="/installation/common/01-flavors.md" %}}

{{% readfile file="/installation/common/02-prerequisites.md" %}}

{{% readfile file="/installation/common/03-prebuilt-binaries.md" %}}

## Package managers

{{% readfile file="/installation/common/homebrew.md" %}}

### MacPorts

[MacPorts] is a free and open source package manager for macOS. This will install the extended flavor of Hugo:

```sh
sudo port install hugo
```

[MacPorts]: https://www.macports.org/

{{% readfile file="/installation/common/04-docker.md" %}}

{{% readfile file="/installation/common/05-build-from-source.md" %}}

## Comparison

||Prebuilt binaries|Package managers|Docker|Build from source
:--|:--:|:--:|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^1]|:heavy_check_mark:|:heavy_check_mark:
Automatic updates?|:x:|:x: [^2]|:x: [^2]|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:

[^1]: Easy if a previous version is still installed.
[^2]: Possible but requires advanced configuration.
