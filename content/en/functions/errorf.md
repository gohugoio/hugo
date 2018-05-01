---
title: errorf
linktitle: errorf
description: Evaluates a format string and logs it to ERROR.
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

`errorf` will evaluate a format string, then output the result to the ERROR log.
This will also cause the build to fail.

```
{{ errorf "Something went horribly wrong! %s" err }}
```
