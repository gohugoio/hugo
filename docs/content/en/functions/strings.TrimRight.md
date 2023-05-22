---
title: strings.TrimRight
description: Returns a slice of a given string with all trailing characters contained in the cutset removed.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["strings.TrimRight CUTSET STRING"]
relatedfuncs: [strings.TrimRight]
---

Given the string `"abba"`, trailing `"a"`'s can be removed a follows:

    {{ strings.TrimRight "a" "abba" }} → "abb"

Numbers can be handled as well:

    {{ strings.TrimRight 12 1221341221 }} → "122134"
