---
title: Logging
description: Enable logging to inspect events while building your site.
categories: [troubleshooting]
keywords: []
menu:
  docs:
    parent: troubleshooting
    weight: 30
weight: 30
toc: true
---

## Command line

Enable console logging with the `--logLevel` command line flag.

Hugo has four logging levels:

error
: Display error messages only.

```sh
hugo --logLevel error
```

warn
: Display warning and error messages.

```sh
hugo --logLevel warn
```

info
: Display information, warning, and error messages.

```sh
hugo --logLevel info
```

debug
: Display debug, information, warning, and error messages.

```sh
hugo --logLevel debug
```

{{% note %}}
If you do not specify a logging level with the `--logLevel` flag, warnings and errors are always displayed.
{{% /note %}}

## Template functions

You can also use template functions to print warnings or errors to the console. These functions are typically used to report data validation errors, missing files, etc.

{{< list-pages-in-section path=/functions/fmt filter=functions_fmt_logging filterType=include >}}
