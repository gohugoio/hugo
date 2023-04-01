---
title: strings.HasSuffix
description: Determine whether a given string ends with the provided trailing suffix string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["strings.HasSuffix STRING SUFFIX"]
relatedfuncs: [hasPrefix]
---

    {{ $pdfPath := "/path/to/some.pdf" }}
    {{ strings.HasSuffix $pdfPath "pdf" }} → true
    {{ strings.HasSuffix $pdfPath "txt" }} → false
