---
title: strings.TrimLeft
description: Returns a slice of a given string with all leading characters contained in the cutset removed.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
signature:
  - strings.TrimLeft CUTSET STRING

---

Given the string `"abba"`, leading `"a"`'s can be removed a follows:

    {{ strings.TrimLeft "a" "abba" }} → "bba"

Numbers can be handled as well:

    {{ strings.TrimLeft 12 1221341221 }} → "341221"
