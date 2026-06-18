---
title: Windows
description: Install Hugo on Windows.
categories: []
keywords: []
weight: 30
---

> [!NOTE]
> Hugo requires Windows 10, Windows Server 2016, or later.

{{% include "/_common/installation/01-editions.md" %}}

{{% include "/_common/installation/02-prerequisites.md" %}}

{{% include "/_common/installation/03-prebuilt-binaries.md" %}}

## Package managers

### Chocolatey

[Chocolatey][] is a free and open-source package manager for Windows. To install the extended edition of Hugo:

```sh
choco install hugo-extended
```

### Scoop

[Scoop][] is a free and open-source package manager for Windows. To install the extended edition of Hugo:

```sh
scoop install hugo-extended
```

### Winget

[Winget][] is Microsoft's official free and open-source package manager for Windows. To install the extended edition of Hugo:

```sh
winget install Hugo.Hugo.Extended
```

To uninstall the extended edition of Hugo:

```sh
winget uninstall --name "Hugo (Extended)"
```

## Build from source

To build Hugo from source you must install:

1. [Git][]
1. [Go][] version {{% current-go-version %}} or later

> [!NOTE]
> The Bash-style `KEY=VALUE cmd` syntax used in the macOS and Linux build-from-source instructions does not work in PowerShell or Command Prompt. Use the code block matching your shell.

### Standard edition

To build and install the standard edition:

PowerShell:

```powershell
$env:CGO_ENABLED=0; go install github.com/gohugoio/hugo@latest
```

Command Prompt:

```bat
set CGO_ENABLED=0
go install github.com/gohugoio/hugo@latest
```

### Deploy edition

{{< new-in v0.159.2 />}}

To build and install the deploy edition:

PowerShell:

```powershell
$env:CGO_ENABLED=0; go install -tags withdeploy github.com/gohugoio/hugo@latest
```

Command Prompt:

```bat
set CGO_ENABLED=0
go install -tags withdeploy github.com/gohugoio/hugo@latest
```

### Extended edition

To build and install the extended edition, first install a C compiler such as [GCC][] or [Clang][] and then run the following command:

PowerShell:

```powershell
$env:CGO_ENABLED=1; go install -tags extended github.com/gohugoio/hugo@latest
```

Command Prompt:

```bat
set CGO_ENABLED=1
go install -tags extended github.com/gohugoio/hugo@latest
```

### Extended/deploy edition

To build and install the extended/deploy edition, first install a C compiler such as [GCC][] or [Clang][] and then run the following command:

PowerShell:

```powershell
$env:CGO_ENABLED=1; go install -tags extended,withdeploy github.com/gohugoio/hugo@latest
```

Command Prompt:

```bat
set CGO_ENABLED=1
go install -tags extended,withdeploy github.com/gohugoio/hugo@latest
```

> [!NOTE]
> See these [detailed instructions][] to install GCC on Windows.

## Comparison

&nbsp;|Prebuilt binaries|Package managers|Build from source
:--|:--:|:--:|:--:
Easy to install?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to upgrade?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Easy to downgrade?|:heavy_check_mark:|:heavy_check_mark: [^2]|:heavy_check_mark:
Automatic updates?|:x:|:x: [^1]|:x:
Latest version available?|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:

[^1]: Possible but requires advanced configuration.
[^2]: Easy if a previous version is still installed.

[Chocolatey]: https://chocolatey.org/
[Clang]: https://clang.llvm.org/
[GCC]: https://gcc.gnu.org/
[Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Go]: https://go.dev/doc/install
[Scoop]: https://scoop.sh/
[Winget]: https://learn.microsoft.com/en-us/windows/package-manager/
[detailed instructions]: https://discourse.gohugo.io/t/41370
