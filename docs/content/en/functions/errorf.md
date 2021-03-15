---
title: errorf and warnf
description: Log ERROR or WARNING from the templates.
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
---

`errorf` or `warnf` will evaluate a format string, then output the result to the ERROR or WARNING log (and only once per error message to avoid flooding the log).

Any ERROR will also cause the build to fail (the `hugo` command will `exit -1`).

Both functions return an empty string, so the messages are only printed to the console.

```
{{ errorf "Failed to handle page %q" .Path }}
```

```
{{ warnf "You should update the shortcodes in %q" .Path }}
```

Note that `errorf` and `warnf` support all the formatting verbs of the [fmt](https://golang.org/pkg/fmt/) package.
