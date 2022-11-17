---
title: strings.HasSuffix
description: Determine whether a given string ends with the provided trailing suffix string.
date: 2019-08-13
publishdate: 2019-08-13
lastmod: 2019-08-13
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["strings.HasSuffix STRING SUFFIX"]
workson: []
hugoversion:
relatedfuncs: [hasPrefix]
deprecated: false
aliases: []
---

    {{ $pdfPath := "/path/to/some.pdf" }}
    {{ strings.HasSuffix $pdfPath "pdf" }} → true
    {{ strings.HasSuffix $pdfPath "txt" }} → false
