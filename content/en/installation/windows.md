---
title: Windows
linkTitle: Windows
description: Install Hugo on Windows.
categories: [installation]
menu:
  docs:
    parent: installation
    weight: 40
toc: true
weight: 40
---
{{% readfile file="/installation/common/01-flavors.md" %}}

{{% readfile file="/installation/common/02-prerequisites.md" %}}

{{% readfile file="/installation/common/03-prebuilt-binaries.md" %}}

## Package managers

### Chocolatey

[Chocolatey] is a free and open source package manager for Windows. This will install the extended flavor of Hugo:

```sh
choco install hugo-extended
```

[Chocolatey]: https://chocolatey.org/

### Scoop

[Scoop] is a free and open source package manager for Windows. This will install the extended flavor of Hugo:

```sh
scoop install hugo-extended
```

[Scoop]: https://scoop.sh/

{{% readfile file="/installation/common/04-docker.md" %}}

{{% readfile file="/installation/common/05-build-from-source.md" %}}

{{% note %}}
When building the extended flavor of Hugo from source on Windows, you will also need to install the [GCC compiler]. See these [detailed instructions].

[detailed instructions]: https://discourse.gohugo.io/t/41370
[GCC compiler]: https://gcc.gnu.org/
{{% /note %}}

## Comparison

||Prebuilt binaries|Package managers|Docker|Build from source
:--|:--:|:--:|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^2]|:heavy_check_mark:|:heavy_check_mark:
Automatic updates?|:x:|:x: [^1]|:x: [^1]|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:

[^1]: Possible but requires advanced configuration.
[^2]: Easy if a previous version is still installed.
