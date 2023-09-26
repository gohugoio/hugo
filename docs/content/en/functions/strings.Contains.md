---
title: strings.Contains
description: Reports whether a string contains a substring.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [string strings substring contains]
signature: ["strings.Contains STRING SUBSTRING"]
relatedfuncs: [strings.ContainsAny]
---

    {{ strings.Contains "Hugo" "go" }} → true

The check is case sensitive: 

    {{ strings.Contains "Hugo" "Go" }} → false
