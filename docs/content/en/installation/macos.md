---
title: macOS
description: Install Hugo on macOS.
categories: []
keywords: []
weight: 10
---

## Editions

{{% include "/_common/installation/01-editions.md" %}}

Unless your specific deployment needs require the extended/deploy edition, we recommend the extended edition.

{{% include "/_common/installation/02-prerequisites.md" %}}

{{% include "/_common/installation/03-prebuilt-binaries.md" %}}

## Package managers

{{% include "/_common/installation/homebrew.md" %}}

### MacPorts

[MacPorts] is a free and open-source package manager for macOS. To install the extended edition of Hugo:

```sh
sudo port install hugo
```

[MacPorts]: https://www.macports.org/

{{% include "/_common/installation/04-build-from-source.md" %}}

## Comparison

&nbsp;|Prebuilt binaries|Package managers|Build from source
:--|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^1]|:heavy_check_mark:
Automatic updates?|:x:|:x: [^2]|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:

[^1]: Easy if a previous version is still installed.
[^2]: Possible but requires advanced configuration.
