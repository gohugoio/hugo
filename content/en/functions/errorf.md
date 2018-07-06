---
title: errorf
linktitle: errorf
description: Log ERROR and fail the build from the templates.
date: 2017-09-30
publishdate: 2017-09-30
lastmod: 2017-09-30
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings, log, error]
signature: ["errorf FORMAT INPUT"]
workson: []
hugoversion:
relatedfuncs: [printf]
deprecated: false
aliases: []
---

`errorf` will evaluate a format string, then output the result to the ERROR log (and only once per error message to avoid flooding the log).

This will also cause the build to fail (the `hugo` command will `exit -1`).

```
{{ errorf "Failed to handle page %q" .Path }}
```

Note that `errorf` supports all the formatting verbs of the [fmt](https://golang.org/pkg/fmt/) package.
