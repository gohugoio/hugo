---
title: Windows
description: Install Hugo on Windows.
categories: []
keywords: []
weight: 30
---

> [!note]
> Hugo v0.121.1 and later require at least Windows 10 or Windows Server 2016.

## Editions

{{% include "/_common/installation/01-editions.md" %}}

Unless your specific deployment needs require the extended/deploy edition, we recommend the extended edition.

{{% include "/_common/installation/02-prerequisites.md" %}}

{{% include "/_common/installation/03-prebuilt-binaries.md" %}}

## Package managers

### Chocolatey

[Chocolatey] is a free and open-source package manager for Windows. To install the extended edition of Hugo:

```sh
choco install hugo-extended
```

### Scoop

[Scoop] is a free and open-source package manager for Windows. To install the extended edition of Hugo:

```sh
scoop install hugo-extended
```

### Winget

[Winget] is Microsoft's official free and open-source package manager for Windows. To install the extended edition of Hugo:

```sh
winget install Hugo.Hugo.Extended
```

To uninstall the extended edition of Hugo:

```sh
winget uninstall --name "Hugo (Extended)"
```

{{% include "/_common/installation/04-build-from-source.md" %}}

> [!note]
> See these [detailed instructions](https://discourse.gohugo.io/t/41370) to install GCC on Windows.

## Comparison

||Prebuilt binaries|Package managers|Build from source
:--|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^2]|:heavy_check_mark:
Automatic updates?|:x:|:x: [^1]|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:

[^1]: Possible but requires advanced configuration.
[^2]: Easy if a previous version is still installed.

[Chocolatey]: https://chocolatey.org/
[Scoop]: https://scoop.sh/
[Winget]: https://learn.microsoft.com/en-us/windows/package-manager/
