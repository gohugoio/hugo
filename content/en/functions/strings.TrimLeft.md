---
title: strings.TrimLeft
description: Returns a slice of a given string with all leading characters contained in the cutset removed.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["strings.TrimLeft CUTSET STRING"]
workson: []
hugoversion:
relatedfuncs: [strings.TrimRight]
deprecated: false
aliases: []
---

Given the string `"abba"`, leading `"a"`'s can be removed a follows:

    {{ strings.TrimLeft "a" "abba" }} → "bba"

Numbers can be handled as well:

    {{ strings.TrimLeft 12 1221341221 }} → "341221"
