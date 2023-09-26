---
title: strings.Contains
description: Reports whether a string contains a substring.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.Contains
  - strings.ContainsAny
  - strings.ContainsNonSpace
  - strings.HasPrefix
  - strings.HasSuffix
  - collections.In
signature:
  - strings.Contains STRING SUBSTRING
---

    {{ strings.Contains "Hugo" "go" }} → true

The check is case sensitive: 

    {{ strings.Contains "Hugo" "Go" }} → false
