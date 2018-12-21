---
title: strings.TrimRight
description: Returns a slice of a given string with all trailing characters contained in the cutset removed.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["strings.TrimRight CUTSET STRING"]
workson: []
hugoversion:
relatedfuncs: [strings.TrimRight]
deprecated: false
aliases: []
---

Given the string `"abba"`, trailing `"a"`'s can be removed a follows:

    {{ strings.TrimRight "a" "abba" }} → "abb"

Numbers can be handled as well:

    {{ strings.TrimRight 12 1221341221 }} → "122134"

